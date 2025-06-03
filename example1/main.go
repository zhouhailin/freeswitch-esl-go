package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/zhouhailin/freeswitch-esl-go/esl"
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
}
func (l *EslConnectionListener) Connected(c *esl.Client) {
}
func (l *EslConnectionListener) Authenticated(authenticated bool, c *esl.Client) {
}
func (l *EslConnectionListener) Disconnected(c *esl.Client) {
}

func main() {
	eventListener := EslEventListener{}
	eslConnectionListener := EslConnectionListener{}
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
	client.AddEventListener(&eventListener)
	client.AddConnectionListener(&eslConnectionListener)
	//client.Close()
	fmt.Println(client)
	time.Sleep(200 * time.Second)

}
