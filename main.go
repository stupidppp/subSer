package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var server *SubServer

func main() {
	server = NewSubServer()
	defer server.Stop()
	//syncTest()
	currencyTest()
	Wait()
}

func Wait() {
	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func syncTest() {

	cli1 := NewClient("tpy", "127.0.0.1")
	cli2 := NewClient("zjj", "127.1.0.2")
	cli3 := NewClient("qdh", "192.22.22.1")

	server.Subscrible("snook", cli1)
	server.Subscrible("snook", cli2)

	server.Subscrible("basket", cli3)
	server.Subscrible("snook", cli3)

	fmt.Println(cli3.subChannels)

	ch := server.GetChannel("snook")
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
	server.UnSubscrible("basket", cli3)
	server.DeleteChannel("basket")
	server.Publish("basket", "rocket failed")
	server.Publish("snook", "奥沙利文 victory 2")
}

func currencyTest() {

	clients := make(chan *Client, 1000)
	for i := 0; i < 1000; i++ {
		cli := NewClient(fmt.Sprintf("%d", i), "127.0.0.1")
		clients <- cli
	}

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				server.Subscrible("snook", <-clients)
			}(&wg)
		} else {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				server.Subscrible("basket", <-clients)
			}(&wg)
		}
	}
	wg.Wait()
	fmt.Println("snook")
	server.GetChannel("snook").PrintClients()
	fmt.Println("basket")
	server.GetChannel("basket").PrintClients()

	server.DeleteChannel("basket")

	cli := server.GetChannel("snook").GetClient("1")
	if cli != nil {
		server.UnSubscrible("snook", cli)
	}

	go server.Publish("snook", "奥沙利文vs塞尔比")
	go server.Publish("basket", "rokect!")
}
