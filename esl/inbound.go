package esl

import (
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

type Client struct {
	gnet.BuiltinEventEngine
	SocketConnection
	gnetClient          *gnet.Client
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
	Level                    Level
}

var options = Options{
	AutoReconnection:         true,
	ReconnectIntervalSeconds: 5,
	MaxReconnectAttempts:     0,
	Level:                    InfoLevel,
}

type ProtocolListener struct{}

func (l ProtocolListener) authResponseReceived(c *Client, response *CommandResponse) {
	c.authenticatorResponded = true
	c.authenticated = response.IsOk()
	c.authenticationResponse = response
	if isDebugEnabled() {
		logging.Debugf("Auth response success=%s, message=[%s]", strconv.FormatBool(c.authenticated), response.GetReplyText())
	}
}

func (l ProtocolListener) eventReceived(c *Client, event *EslEvent) {
	if isDebugEnabled() {
		logging.Debugf("Event received %s", event.ToString())
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
					logging.Errorf("%d Error caught notifying listener of job result %s: %v", i, event.ToString(), err)
				}
			}
		} else {
			for i, listener := range c.eventListeners {
				err := listener.EventReceived(event)
				if err != nil {
					logging.Errorf("%d Error caught notifying listener of event %s: %v", i, event.ToString(), err)
				}
			}
		}
	}()
}

func (l ProtocolListener) disconnected(c *Client) {
	if isInfoEnabled() {
		logging.Infof("Disconnected ..")
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

// OnOpen is called when a new connection is opened.
func (client *Client) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	if isTraceEnabled() {
		logging.Debugf("OnOpen: %v", c.RemoteAddr())
	}
	return nil, gnet.None
}

// OnTraffic is called when data is available on the connection.
// It decodes ESL messages and dispatches them to the appropriate handler.
func (client *Client) OnTraffic(c gnet.Conn) (action gnet.Action) {
	for {
		m := EslMessage{
			headers:       make(map[Name]string),
			body:          []string{},
			contentLength: 0,
		}
		err := decode(c, &m)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				// Not enough data for a complete message, wait for more
				return gnet.None
			}
			logging.Errorf("Decode error: %v", err)
			return gnet.Close
		}
		// Process message in a goroutine to avoid blocking the event loop
		go messageReceived(client, &m)
	}
}

// OnClose is called when a connection is closed.
func (client *Client) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	// Only handle if this is the current connection
	if client.conn != c {
		return gnet.Close
	}
	client.active = false
	logging.Infof("[%v] connection closed", c.RemoteAddr())
	close(client.msg)
	// Notify connection is disconnect
	if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
		go func() {
			for _, listener := range client.connectionListeners {
				listener.Disconnected(client)
			}
		}()
	}
	// reconnect
	client.canReconnect()
	return gnet.Close
}

func (client *Client) Connect() error {
	if client.CanSend() {
		if isInfoEnabled() {
			logging.Infof("Client is connected, will close first.")
		}
		_, err := client.Close()
		if err != nil {
			return err
		}
		time.Sleep(250 * time.Millisecond)
	}
	// Create gnet client if not exists
	if client.gnetClient == nil {
		gc, err := gnet.NewClient(client, gnet.WithLogLevel(options.Level))
		if err != nil {
			if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
				go func() {
					for _, listener := range client.connectionListeners {
						listener.ConnectFailure(client)
					}
				}()
			}
			client.canReconnect()
			return err
		}
		err = gc.Start()
		if err != nil {
			if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
				go func() {
					for _, listener := range client.connectionListeners {
						listener.ConnectFailure(client)
					}
				}()
			}
			client.canReconnect()
			return err
		}
		client.gnetClient = gc
	}
	// Dial connection
	conn, err := client.gnetClient.Dial(client.Network, client.Address)
	if err != nil {
		if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
			go func() {
				for _, listener := range client.connectionListeners {
					listener.ConnectFailure(client)
				}
			}()
		}
		client.canReconnect()
		return err
	}
	// Setup SocketConnection fields
	client.conn = conn
	client.msg = make(chan *EslMessage)
	client.authenticationResponse = nil
	client.authenticatorResponded = false
	client.authenticated = false
	client.rudeRejection = false
	client.active = true
	client.listener = ProtocolListener{}

	if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
		go func() {
			for _, listener := range client.connectionListeners {
				listener.Connected(client)
			}
		}()
	}

	for !client.rudeRejection && !client.authenticatorResponded {
		time.Sleep(250 * time.Millisecond)
	}

	if client.connectionListeners != nil && len(client.connectionListeners) > 0 {
		go func() {
			for _, listener := range client.connectionListeners {
				listener.Authenticated(client.authenticated, client)
			}
		}()
	}
	if client.rudeRejection {
		return errors.New("client is rejected by acl")
	} else if !client.authenticated {
		return errors.New("Authentication failed: " + client.authenticationResponse.GetReplyText())
	}
	return nil
}

func (client *Client) canReconnect() {
	if options.AutoReconnection && options.ReconnectIntervalSeconds > 0 {
		time.AfterFunc(time.Duration(options.ReconnectIntervalSeconds)*time.Second, func() {
			logging.Infof("Reconnecting ...")
			err := client.Connect()
			if err != nil {
				logging.Errorf("Reconnection failure, cause %v", err)
			}
		})
	}
}
