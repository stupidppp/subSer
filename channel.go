package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Channel struct {
	name     string             // 频道名
	clients  map[string]*Client // 订阅频道的客户
	msg      chan string        // 发送消息队列
	quitChan chan bool          // 退出提示
	wg       sync.WaitGroup     // 监听发送队列的服务
	lock     sync.RWMutex
	exit     uint32 // 结束标记
}

func (ch *Channel) ExitAndSet() bool {
	return atomic.CompareAndSwapUint32(&ch.exit, 0, 1)
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:     name,
		clients:  make(map[string]*Client, 32),
		msg:      make(chan string, 100),
		quitChan: make(chan bool, 1),
		exit:     uint32(0),
	}
}

func (ch *Channel) Start() {
	go ch.handleNotify()
}

// thread safe
func (ch *Channel) Close() {
	if ch.ExitAndSet() {
		if ch.Length() <= 0 {
			close(ch.quitChan)
		} else {
			close(ch.msg) // 继续发送队列剩余的消息
		}
		ch.wg.Wait()
	}
}

func (ch *Channel) Notify(msg string) {
	// may be close channel panic
	defer func() {
		recover()
	}()

	ch.msg <- msg
}

func (ch *Channel) handleNotify() {
	ch.wg.Add(1)
	defer ch.wg.Done()

	for {
		select {
		case <-ch.quitChan:
			return

		case msg, ok := <-ch.msg:
			if !ok {
				return
			}
			ch.lock.RLock()
			for n, _ := range ch.clients {
				fmt.Printf("channel[%s]to client[%s] notify[%s] \n", ch.name, n, msg)
			}
			ch.lock.RUnlock()
		}
	}
}

func (ch *Channel) AddClient(c *Client) {
	if c == nil || c.uid == "" {
		return
	}
	ch.lock.RLock()
	_, found := ch.clients[c.uid]
	ch.lock.RUnlock()

	if !found {
		ch.lock.Lock()
		if _, found := ch.clients[c.uid]; !found {
			ch.clients[c.uid] = c
		}
		ch.lock.Unlock()
	}
}

func (ch *Channel) DeleteClient(c *Client) {
	if c == nil || c.uid == "" {
		return
	}
	ch.lock.Lock()
	if c1, found := ch.clients[c.uid]; found && c1 == c {
		delete(ch.clients, c.uid)
	}
	ch.lock.Unlock()
}

func (ch *Channel) GetClient(uname string) *Client {
	if uname == "" {
		return nil
	}
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	c, found := ch.clients[uname]
	if found {
		return c
	}
	return nil
}

func (ch *Channel) Length() int {
	ch.lock.RLock()
	defer ch.lock.RUnlock()
	return len(ch.clients)
}

func (ch *Channel) PrintClients() {
	ch.lock.RLock()
	for n, _ := range ch.clients {
		fmt.Println("client:", n)
	}
	ch.lock.RUnlock()
}
