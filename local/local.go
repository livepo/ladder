package local

import (
    "fmt"
    "ladder/cipher"
    "net"
    "bufio"
    "io"
)


type Local struct {
    LocalAddr string
    RemoteAddr string
}

func NewLocal(localaddr, remoteaddr string) *Local {
    return &Local{
        LocalAddr: localaddr,
        RemoteAddr: remoteaddr,
    }
}


func checkError(err error) {
    if err != nil && err != io.EOF {
        panic(err)
    }
}

func (l *Local) Serve() {
    ln, err := net.Listen("tcp", l.LocalAddr)
    checkError(err)

    for {
        client, err := ln.Accept()
        checkError(err)
        go l.handleClient(client)
    }
}


func (l *Local) handleClient(client net.Conn) {
    browserMsg := make(chan []byte)
    go l.readClient(client, browserMsg)  // 读客户端数据, 加密，生产者
    go l.transportClient(client, browserMsg)  // 转发客户端数据, 回写客户端, 消费者
}


func (l *Local) readClient(client net.Conn, browserMsg chan []byte) {
    defer client.Close()
    defer close(browserMsg)

    reader := bufio.NewReader(client)
    buf := make([]byte, 1024)
    for {
        nsize, err := reader.Read(buf)
        if err != nil && err != io.EOF {
            return
        }
        if nsize > 0 {
            cipher.Encode(buf[:nsize])
            browserMsg <- buf[:nsize]
        }
    }
}


func (l *Local) transportClient(client net.Conn, browserMsg chan []byte) {
    remote, err := net.Dial("tcp", l.RemoteAddr)
    checkError(err)
    go l.remoteToClient(client, remote)
    for message := range browserMsg {
        nsize, err := remote.Write(message)
        if err != nil {
            fmt.Println("write to remote err", err)
            return
        }
        if nsize < len(message) {
            fmt.Println("write to remote nsize error")
        }
    }
}


func (l *Local) remoteToClient(client net.Conn, remote net.Conn) {
    buf := make([]byte, 1024)
    reader := bufio.NewReader(remote)
    for {
        nsize, err := reader.Read(buf)
        if err != nil && err != io.EOF {
            fmt.Println("read remote error", err)
            return
        }
        if nsize > 0 {
            cipher.Decode(buf[:nsize])
            client.Write(buf[:nsize])
        }
    }
}

