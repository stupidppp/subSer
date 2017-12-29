package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := NewSubServer()
	defer server.Stop()
	cli1 := NewClient("tpy", "127.0.0.1")
	cli2 := NewClient("zjj", "127.1.0.2")
	cli3 := NewClient("qdh", "192.22.22.1")

	server.Subscrible("snook", cli1)
	server.Subscrible("snook", cli2)

	server.Subscrible("basket", cli3)
	server.Subscrible("snook", cli3)

	fmt.Println(cli3.subChannels)

	/*ch := server.GetChannel("snook")
	fmt.Println("snook-------------")
	ch.PrintClients()

	ch1 := server.GetChannel("basket")
	fmt.Println("basket------------")
	ch1.PrintClients()

	server.Publish("snook", "奥沙利文vs塞尔比")
	server.Publish("basket", "rokect!")

	fmt.Println("############### delete user")
	ch.DeleteClient(cli3)
	server.Publish("snook", "奥沙利文 victory")

	fmt.Println("##############  delete channel")
	//server.UnSubscrible("basket", cli3)
	//server.DeleteChannel("basket")
	//	server.Publish("basket", "rocket failed")
	//	server.Publish("snook", "奥沙利文 victory 2")
	*/
	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

}
