package server

import (
    "fmt"
    "ladder/cipher"
    "net"
    "bufio"
    "io"
    "encoding/binary"
)


type Server struct {
    Addr string
}


func NewServer(addr string) *Server {
    return &Server{
        Addr: addr,
    }
}


func checkError(err error) {
    if err != nil && err != io.EOF {
        panic(err)
    }
}

func (s *Server) Serve() {
    ln, err := net.Listen("tcp", s.Addr)
    checkError(err)
    for {
        proxy, err := ln.Accept()
        checkError(err)
        go s.HandleProxy(proxy)
    }
}


func (s *Server) HandleProxy(proxy net.Conn) {

    reader := bufio.NewReader(proxy)
    buf := make([]byte, 1024)
    nsize, err := reader.Read(buf)
    if err != nil && err != io.EOF {
        return
    }
    if nsize > 0 {
        cipher.Decode(buf[:nsize])
        if buf[0] != 0x05 {
            return
        }
        authorized := []byte{0x05, 0x00}
        cipher.Encode(authorized)
        proxy.Write(authorized)
    }
    nsize, err = reader.Read(buf)
    if err != nil && err != io.EOF {
        return
    }
    var realip []byte
    if nsize > 0 {
        cipher.Decode(buf[:nsize])
        switch buf[3] {
        case 0x01:
            // IP V4 address: X'01'
            realip = buf[4:4+net.IPv4len]
        case 0x03:
            // DOMAINNAME: X'03'
            ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:nsize-2]))
            if err != nil {
                return
            }
            realip = ipAddr.IP
        case 0x04:
            // IP V6 address: X'04'
            realip = buf[4:4+net.IPv6len]
        default:
            return
        }
        host := net.IP(realip).String()
        port := int(binary.BigEndian.Uint16(buf[nsize-2:]))
        remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
        if err != nil {
            fmt.Println("connect remote error", err)
            return
        }
        connected := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
        cipher.Encode(connected)
        proxy.Write(connected)
        go Transport(remote, proxy, cipher.Decode)
        go Transport(proxy, remote, cipher.Encode)
    }
}


type CiphAction func([]byte)


func Transport(dst net.Conn, src net.Conn, ca CiphAction) {
    reader := bufio.NewReader(src)
    buf := make([]byte, 1024)
    for {
        nsize, err := reader.Read(buf)
        if err != nil && err != io.EOF {
            fmt.Println("read err", err)
            return
        }
        if nsize > 0 {
            ca(buf)
            dst.Write(buf[:nsize])
        } else {
            return
        }
    }
}
