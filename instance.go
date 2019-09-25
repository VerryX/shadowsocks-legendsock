package main

import (
	"log"
	"net"

	"github.com/shadowsocks-server/shadowsocks-legendsock/core"
)

// Instance struct
type Instance struct {
	UserID    int
	Port      int
	Method    string
	Password  string
	Bandwidth *Bandwidth
	Started   bool
	TCPSocket net.Listener
	UDPSocket net.PacketConn
}

// Start instance
func (instance *Instance) Start() {
	cipher, err := core.PickCipher(instance.Method, make([]byte, 0), instance.Password)
	if err != nil {
		log.Println(err)
		return
	}

	instance.Started = true

	go tcpRemote(instance, cipher.StreamConn)
	go udpRemote(instance, cipher.PacketConn)
}

// Stop instance
func (instance *Instance) Stop() {
	instance.Started = false

	if instance.TCPSocket != nil {
		instance.TCPSocket.Close()
	}

	if instance.UDPSocket != nil {
		instance.UDPSocket.Close()
	}
}

func newInstance(id int, port int, method, password string) *Instance {
	instance := Instance{}
	instance.UserID = id
	instance.Port = port
	instance.Method = method
	instance.Password = password
	instance.Bandwidth = newBandwidth()
	instance.Started = false

	return &instance
}
