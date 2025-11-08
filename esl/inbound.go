package esl

import (
	"context"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/netpoll"
	"net"
	"strconv"
	"time"
)

type Client struct {
	SocketConnection
	Network             string
	Address             string
	Password            string
	TimeoutSeconds      int
	reconnectAttempts   int
	eventListeners      []IEslEventListener
	connectionListeners []IEslConnectionListener
}

type Options struct {
	AutoReconnection         bool
	ReconnectIntervalSeconds int
	MaxReconnectAttempts     int
	Level                    logger.Level
}

var options = Options{
	AutoReconnection:         true,
	ReconnectIntervalSeconds: 5,
	MaxReconnectAttempts:     0,
	Level:                    logger.LevelInfo,
}

type ProtocolListener struct {
}

func (l ProtocolListener) authResponseReceived(c *Client, response *CommandResponse) {
	c.authenticatorResponded = true
	c.authenticated = response.IsOk()
	c.authenticationResponse = response
	if isDebugEnabled() {
		logger.Debug("Auth response success=" + strconv.FormatBool(c.authenticated) + ", message=[" + response.GetReplyText() + "]")
	}
}

func (l ProtocolListener) eventReceived(c *Client, event *EslEvent) {
	// log.debug( "Event received [{}]", event );
	if isDebugEnabled() {
		logger.Debugf("Event received %s\n", event.ToString())
	}
	if c.eventListeners == nil || len(c.eventListeners) == 0 {
		return
	}

	/*
	 *  Notify listeners in a different thread in order to:
	 *    - not to block the IO threads with potentially long-running listeners
	 *    - generally be defensive running other people's code
	 *  Use a different worker thread pool for async job results than for event driven
	 *  events to keep the latency as low as possible.
	 */
	go func() {
		if event.GetEventName() == "BACKGROUND_JOB" {
			for i, listener := range c.eventListeners {
				err := listener.BackgroundJobResultReceived(event)
				if err != nil {
					logger.Errorf("%d Error caught notifying listener of job result %s\n", i, event.ToString(), err)
				}
			}
		} else {
			for i, listener := range c.eventListeners {
				err := listener.EventReceived(event)
				if err != nil {
					logger.Errorf("%d Error caught notifying listener of event %s\n", i, event.ToString(), err)
				}
			}
		}
	}()
}

func (l ProtocolListener) disconnected(c *Client) {
	if isInfoEnabled() {
		logger.Info("Disconnected ..")
	}
}

// NewClient - Will initiate new client that will establish connection and attempt to authenticate
// @Param host
func NewClient(host string, port uint, password string, timeoutSeconds int, newOptions *Options) *Client {
	if newOptions != nil {
		options = *newOptions
	}
	return &Client{
		Network:             "tcp",
		Address:             net.JoinHostPort(host, strconv.Itoa(int(port))),
		Password:            password,
		TimeoutSeconds:      timeoutSeconds,
		reconnectAttempts:   0,
		eventListeners:      nil,
		connectionListeners: nil,
	}
}

func (client *Client) AddEventListener(listener IEslEventListener) {
	if client.eventListeners == nil {
		client.eventListeners = *new([]IEslEventListener)
	}
	client.eventListeners = append(client.eventListeners, listener)
}

func (client *Client) AddConnectionListener(listener IEslConnectionListener) {
	if client.connectionListeners == nil {
		client.connectionListeners = *new([]IEslConnectionListener)
	}
	client.connectionListeners = append(client.connectionListeners, listener)
}

func (client *Client) Connect() error {
	if client.CanSend() {
		if isInfoEnabled() {
			logger.Info("Client is connected, will close first.")
		}
		_, err := client.Close()
		if err != nil {
			return err
		}
	}
	// use default
	connection, err := netpoll.DialConnection(client.Network, client.Address, time.Duration(client.TimeoutSeconds)*time.Second)
	if err != nil {
		go func() {
			client.canReconnect()
			if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
				for _, listener := range client.connectionListeners {
					listener.ConnectFailure(client)
				}
			}
		}()
		return err
	}
	if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
		go func() {
			for _, listener := range client.connectionListeners {
				listener.Connected(client)
			}
		}()
	}
	client.SocketConnection = SocketConnection{
		Connection:             connection,
		msg:                    make(chan *EslMessage),
		authenticationResponse: nil,
		authenticatorResponded: false,
		authenticated:          false,
		listener:               ProtocolListener{},
	}
	err = connection.SetOnRequest(func(ctx context.Context, connection netpoll.Connection) error {
		var err error
		if isTraceEnabled() {
			logger.Trace("Connect SetOnRequest .....")
		}
		m := EslMessage{
			headers:       make(map[Name]string),
			body:          *new([]string),
			contentLength: 0,
		}
		err = decode(connection.Reader(), &m)
		if err != nil {
			return err
		}
		return messageReceived(client, &m)
	})
	if err != nil {
		return err
	}

	err = connection.AddCloseCallback(func(connection netpoll.Connection) error {
		logger.Infof("[%v] connection closed\n", connection.RemoteAddr())
		client.authenticationResponse = nil
		client.authenticatorResponded = false
		client.authenticated = false
		close(client.msg)
		// Notify connection is disconnect
		if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
			go func() {
				client.canReconnect()
				for _, listener := range client.connectionListeners {
					listener.Disconnected(client)
				}
			}()
		}
		return nil
	})

	for !client.authenticatorResponded {
		time.Sleep(250 * time.Millisecond)
	}

	if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
		go func() {
			for _, listener := range client.connectionListeners {
				listener.Authenticated(client.authenticated, client)
			}
		}()
	}
	if !client.authenticated {
		return errors.New("Authentication failed: " + client.authenticationResponse.GetReplyText())
	}
	return err
}

func (client *Client) canReconnect() {
	if options.AutoReconnection && options.ReconnectIntervalSeconds > 0 {
		time.AfterFunc(time.Duration(options.ReconnectIntervalSeconds)*time.Second, func() {
			err := client.Connect()
			if err != nil {
				logger.Error("Reconnection failure, cause ", err)
			}
		})
	}
}
