# FreeSWITCH Event Socket Library Go

[![License](https://img.shields.io/github/license/zhouhailin/freeswitch-esl-go)](https://github.com/zhouhailin/freeswitch-esl-go/blob/master/LICENSE)
[![Go 1.15 Ready](https://img.shields.io/badge/Go%201.15-Ready-green.svg?style=flat)]()

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

## Netpoll

[Netpoll][Netpoll] is a high-performance non-blocking I/O networking framework, which focused on RPC scenarios,
developed by [ByteDance][ByteDance].

[Netpoll]: https://github.com/cloudwego/netpoll

## License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) Copyright (C) Apache Software Foundation
