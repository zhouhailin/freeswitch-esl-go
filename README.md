# FreeSWITCH Event Socket Library Go

[![License](https://img.shields.io/github/license/zhouhailin/freeswitch-esl-go)](https://github.com/zhouhailin/freeswitch-esl-go/blob/master/LICENSE)
[![Go 1.15 Ready](https://img.shields.io/badge/Go%201.15-Ready-green.svg?style=flat)]()

## Features

* **Already**
    - Inbound Client
    - Linux, macOS (operating system)

* **Unsupported**
    - Windows (operating system)

## Quick Start

```
go get -u github.com/zhouhailin/freeswitch-esl-go
```

```go
package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/zhouhailin/freeswitch-esl-go/esl"
	"time"
)

func main() {
	// eventListener := EslEventListener{}
	// eslConnectionListener := EslConnectionListener{}
	client := esl.NewClient("127.0.0.1", 8021, "ClueCon", 5, &esl.Options{
		Level: logger.LevelTrace,
	})
	fmt.Println(client)
	//client.Connect()
	err := client.Connect()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	subscriptions, err := client.SetEventSubscriptions("plain", "ALL")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Println(subscriptions)
	// client.AddEventListener(&eventListener)
	// client.AddConnectionListener(&eslConnectionListener)
	// client.Close()
	fmt.Println(client)
	time.Sleep(200 * time.Second)
}
```

```go
package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/zhouhailin/freeswitch-esl-go/esl"
	"os"
	"strconv"
	"time"
)

type EslEventListener struct {
}

func (l *EslEventListener) EventReceived(event *esl.EslEvent) error {
	fmt.Println("######## eventReceived : " + event.ToString())
	return nil
}

func (l *EslEventListener) BackgroundJobResultReceived(event *esl.EslEvent) error {
	fmt.Println("######## backgroundJobResultReceived : " + event.ToString())
	return nil
}

type EslConnectionListener struct {
}

func (l *EslConnectionListener) ConnectFailure(c *esl.Client) {
	fmt.Println("ConnectFailure")
}
func (l *EslConnectionListener) Connected(client *esl.Client) {
	fmt.Println("Connected")
}
func (l *EslConnectionListener) Authenticated(authenticated bool, client *esl.Client) {
	fmt.Println("Authenticated : " + strconv.FormatBool(authenticated))
	if authenticated {
		subscriptions, err := client.SetEventSubscriptions("plain", "ALL")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Println(subscriptions)
	}
}
func (l *EslConnectionListener) Disconnected(c *esl.Client) {
	fmt.Println("Disconnected")
}

func main() {
	eventListener := EslEventListener{}
	eslConnectionListener := EslConnectionListener{}
	env, b := os.LookupEnv("PATH")
	println(env, b)
	client := esl.NewClient("127.0.0.1", 8021, "ClueCon", 5, &esl.Options{
		AutoReconnection:         true,
		ReconnectIntervalSeconds: 5,
		Level:                    logger.LevelDebug,
	})
	fmt.Println(client)
	client.AddEventListener(&eventListener)
	client.AddConnectionListener(&eslConnectionListener)

	err := client.Connect()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	//client.Close()
	fmt.Println(client)
	time.Sleep(200 * time.Second)
}
```

## Netpoll

[Netpoll][Netpoll] is a high-performance non-blocking I/O networking framework, which focused on RPC scenarios,
developed by [ByteDance][ByteDance].

[Netpoll]: https://github.com/cloudwego/netpoll

## License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) Copyright (C) Apache Software Foundation
