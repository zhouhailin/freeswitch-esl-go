# FreeSWITCH Event Socket Library Go

[![License](https://img.shields.io/github/license/zhouhailin/freeswitch-esl-go)](https://github.com/zhouhailin/freeswitch-esl-go/blob/master/LICENSE)
[![Go 1.25 Ready](https://img.shields.io/badge/Go%201.25-Ready-green.svg?style=flat)]()

## Features

* **Already**
    - Inbound Client
    - Linux, macOS (operating system)
    - Windows (development and testing only)

* **Unsupported**
    - Windows (production environment)

## Quick Start

```
go get -u github.com/zhouhailin/freeswitch-esl-go
```

```go
package main

import (
	"fmt"
	"github.com/zhouhailin/freeswitch-esl-go/esl"
	"os"
	"strconv"
	"time"
)

type EslEventListener struct {
}

func printlnEslEvent(event *esl.EslEvent) {
	fmt.Println("###### messageHeaders ")
	for name, value := range *event.GetMessageHeaders() {
		fmt.Printf("name : [%s], value : [%s]\n", name, value)
	}
	fmt.Println("###### eventHeaders ")
	for name, value := range *event.GetEventHeaders() {
		fmt.Printf("name : [%s], value : [%s]\n", name, value)
	}
	fmt.Println("###### eventBody ")
	for index, line := range *event.GetEventBodyLines() {
		fmt.Printf(" [%d] : [%s]\n", index, line)
	}
}

func (l *EslEventListener) EventReceived(event *esl.EslEvent) error {
	fmt.Println("######## eventReceived : " + event.ToString())
	printlnEslEvent(event)
	return nil
}

func (l *EslEventListener) BackgroundJobResultReceived(event *esl.EslEvent) error {
	fmt.Println("######## backgroundJobResultReceived : " + event.ToString())
	printlnEslEvent(event)
	return nil
}

type EslConnectionListener struct {
}

func (l *EslConnectionListener) ConnectFailure(c *esl.Client) {
	fmt.Println("ConnectFailure")
}
func (l *EslConnectionListener) Connected(c *esl.Client) {
	fmt.Println("Connected")
}
func (l *EslConnectionListener) Authenticated(authenticated bool, c *esl.Client) {
	fmt.Println("Authenticated : " + strconv.FormatBool(authenticated))
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
		Level: esl.DebugLevel,
	})
	fmt.Println(client)
	client.AddEventListener(&eventListener)
	client.AddConnectionListener(&eslConnectionListener)
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
	fmt.Println(client)
	time.Sleep(200 * time.Second)
}
```

```go
package main

import (
	"fmt"
	"github.com/zhouhailin/freeswitch-esl-go/esl"
	"os"
	"strconv"
	"time"
)

type EslEventListener struct {
}

func printlnEslEvent(event *esl.EslEvent) {
	fmt.Println("###### messageHeaders ")
	for name, value := range *event.GetMessageHeaders() {
		fmt.Printf("name : [%s], value : [%s]\n", name, value)
	}
	fmt.Println("###### eventHeaders ")
	for name, value := range *event.GetEventHeaders() {
		fmt.Printf("name : [%s], value : [%s]\n", name, value)
	}
	fmt.Println("###### eventBody ")
	for index, line := range *event.GetEventBodyLines() {
		fmt.Printf(" [%d] : [%s]\n", index, line)
	}
}

func (l *EslEventListener) EventReceived(event *esl.EslEvent) error {
	fmt.Println("######## eventReceived : " + event.ToString())
	printlnEslEvent(event)
	return nil
}

func (l *EslEventListener) BackgroundJobResultReceived(event *esl.EslEvent) error {
	fmt.Println("######## backgroundJobResultReceived : " + event.ToString())
	printlnEslEvent(event)
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
		ReconnectIntervalSeconds: 3,
		MaxReconnectAttempts:     100,
		Level:                    esl.DebugLevel,
	})
	fmt.Println(client)
	client.AddEventListener(&eventListener)
	client.AddConnectionListener(&eslConnectionListener)
	err := client.Connect()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(client)
	time.Sleep(200 * time.Second)
}
```

## gnet

[gnet][gnet] is a high-performance, lightweight, non-blocking, event-driven networking framework written in pure Go,
developed by [Andy Pan (panjf2000)][panjf2000].

[gnet]: https://github.com/panjf2000/gnet
[panjf2000]: https://github.com/panjf2000

## License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) Copyright (C) Apache Software Foundation
