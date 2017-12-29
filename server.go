package main

import (
	"sync"
)

type SubServer struct {
	channels map[string]*Channel
	lock     sync.RWMutex
}

func NewSubServer() *SubServer {
	return &SubServer{
		channels: make(map[string]*Channel, 100),
	}
}

func (this *SubServer) Stop() {
	this.lock.Lock()
	for _, ch := range this.channels {
		ch.Close()
	}
	this.lock.Unlock()
}

func (this *SubServer) DeleteChannel(chName string) {
	if ch := this.GetChannel(chName); ch == nil {
		return
	}

	this.lock.Lock()
	if ch, ok := this.channels[chName]; ok {
		delete(this.channels, chName)
		ch.Close()
	}
	this.lock.Unlock()
}

func (this *SubServer) GetChannel(chName string) *Channel {
	this.lock.RLock()
	defer this.lock.RUnlock()

	ch, ok := this.channels[chName]
	if ok {
		return ch
	}
	return nil
}

// return 之前储存的channel, 新增的channel
func (this *SubServer) AddChannel(chName string) (pre *Channel, cur *Channel) {
	if ch := this.GetChannel(chName); ch != nil {
		return ch, nil
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	if ch, ok := this.channels[chName]; ok {
		return ch, nil
	}
	channel := NewChannel(chName)
	this.channels[chName] = channel
	return nil, channel
}

func (this *SubServer) Subscrible(chName string, cli *Client) {
	// 0 用户重复订阅
	if cli.ContainChannel(chName) {
		return
	}
	cli.AddChannel(chName)

	pre, cur := this.AddChannel(chName)
	if cur != nil {
		cur.AddClient(cli)
		cur.Start()
		return
	}
	pre.AddClient(cli)
	return
}

func (this *SubServer) UnSubscrible(chName string, cli *Client) {
	if cli.ContainChannel(chName) {
		cli.DeleteChannel(chName)
	}

	ch := this.GetChannel(chName)
	if ch != nil {
		ch.DeleteClient(cli)

		if ch.Length() == 0 {
			this.DeleteChannel(chName)
		}
	}
}

func (this *SubServer) Publish(chName string, msg string) {
	ch := this.GetChannel(chName)
	if ch == nil {
		return
	}
	ch.Notify(msg)
}
