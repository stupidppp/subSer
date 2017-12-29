package main

import (
	"sync"
)

type Client struct {
	uid         string          // 用户名
	ip          string          // 连接
	subChannels map[string]bool // 订阅的频道
	lock        sync.RWMutex
}

func NewClient(uid, ip string) *Client {
	return &Client{
		uid:         uid,
		ip:          ip,
		subChannels: make(map[string]bool)}
}

func (c *Client) AddChannel(name string) {
	c.lock.RLock()
	_, found := c.subChannels[name]
	c.lock.RUnlock()
	if !found {
		c.lock.Lock()
		if _, found := c.subChannels[name]; !found {
			c.subChannels[name] = true
		}
		c.lock.Unlock()
	}
}

func (c *Client) DeleteChannel(name string) {
	c.lock.Lock()
	if _, found := c.subChannels[name]; found {
		delete(c.subChannels, name)
	}
	c.lock.Unlock()
}

func (c *Client) ContainChannel(name string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.subChannels[name]
	return ok
}
