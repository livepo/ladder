package local

import (
	"fmt"
	"gosocks/cipher"
	"gosocks/tunnel"
	"net"
)

type Local struct {
	ListenAddr string
	RemoteAddr string
}

func NewLocal(listenaddr, remoteaddr string) *Local {
	return &Local{
		ListenAddr: listenaddr,
		RemoteAddr: remoteaddr,
	}
}

func (l *Local) Serve() {
	ln, err := net.Listen("tcp", l.ListenAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
		}
		go l.HandleClient(client)
	}
}

func (l *Local) HandleClient(client net.Conn) {
	remote, err := net.Dial("tcp", l.RemoteAddr)
	if err != nil {
		return
	}
	go tunnel.Transport(remote, client, cipher.Encode)
	go tunnel.Transport(client, remote, cipher.Decode)
}
