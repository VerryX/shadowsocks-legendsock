package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/shadowsocks-server/shadowsocks-legendsock/socks"
)

func tcpRemote(instance *Instance, cipher func(net.Conn) net.Conn) {
	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", instance.Port))
	if err != nil {
		log.Printf("Failed to listen TCP on %d: %v", instance.Port, err)
		return
	}
	defer socket.Close()
	instance.TCPSocket = socket

	for instance.Started {
		client, err := socket.Accept()
		if err != nil {
			continue
		}

		go tcpHandle(instance, client, cipher)
	}
}

func tcpHandle(instance *Instance, client net.Conn, cipher func(net.Conn) net.Conn) {
	defer client.Close()
	client.(*net.TCPConn).SetKeepAlive(true)
	client = cipher(client)

	target, err := socks.ReadAddr(client)
	if err != nil {
		return
	}

	remote, err := net.Dial("tcp", target.String())
	if err != nil {
		return
	}
	defer remote.Close()

	tcpRelay(instance, client, remote)
}

func tcpRelay(instance *Instance, left, right net.Conn) {
	go func() {
		size, _ := io.CopyBuffer(right, left, make([]byte, 4096))
		instance.Bandwidth.IncreaseUpload(uint64(size))
		right.SetDeadline(time.Now())
		left.SetDeadline(time.Now())
	}()

	size, _ := io.CopyBuffer(left, right, make([]byte, 4096))
	instance.Bandwidth.IncreaseDownload(uint64(size))
	right.SetDeadline(time.Now())
	left.SetDeadline(time.Now())
}
