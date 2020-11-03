package server


import (
    "fmt"
    "net"
    "gosocks/tunnel"
    "gosocks/cipher"
    "bufio"
    "encoding/binary"
    "errors"
)


type Server struct {
    ListenAddr string
}


func NewServer(listenaddr string) *Server {
    return &Server {
        ListenAddr: listenaddr,
    }
}


func (s *Server) Serve() {
    ln, err := net.Listen("tcp", s.ListenAddr)
    if err != nil {
        fmt.Println("listen error", err)
    }
    for {
        client, err := ln.Accept()
        if err != nil {
            fmt.Println("accept error", err)
        }
        go s.HandleClient(client)
    }
}


func (s *Server) HandleClient(client net.Conn) {
    reader := bufio.NewReader(client)
    buf := make([]byte, 64)
    nsize, err := reader.Read(buf)
    if err != nil {
        fmt.Println("socks5 parse error", err)
        return
    }
    cipher.Decode(buf[:nsize])
    if nsize == 0 || buf[0] != 0x05 {
        fmt.Println("socks5 parse error", err)
        return
    }
    resp1 := []byte{0x05, 0x00}
    cipher.Encode(resp1)
    nsize, err = client.Write(resp1)
    if err != nil {
        fmt.Println("resp socks5 ack error", err)
        return
    }

    nsize, err = reader.Read(buf)
    if err != nil {
        fmt.Println("socks5 parse error", err)
        return
    }
    cipher.Decode(buf)
    remoteAddr, err := s.parseRemote(buf[:nsize])
    if err != nil {
        fmt.Println("parse remote addr error", err)
        return
    }
    remote, err := net.Dial("tcp", remoteAddr)
    if err != nil {
        fmt.Println("connect remote error", err)
        return
    }
    resp2 := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
    cipher.Encode(resp2)
    nsize, err = client.Write(resp2)
    if err != nil {
        fmt.Println("resp socks5 ack error", err)
        return
    }
    go tunnel.Transport(remote, client, cipher.Decode)
    go tunnel.Transport(client, remote, cipher.Encode)
}


func (s *Server) parseRemote(buf []byte) (string, error) {
    var dIP, dPort []byte
    switch buf[3] {
    case 0x01:
        //	IP V4 address: X'01'
        dIP = buf[4 : 4+net.IPv4len]
    case 0x03:
        //	DOMAINNAME: X'03'
        ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:len(buf)-2]))
        if err != nil {
            return "", err
        }
        dIP = ipAddr.IP
    case 0x04:
        //	IP V6 address: X'04'
        dIP = buf[4 : 4+net.IPv6len]
    default:
        return "", errors.New("parse remote addr error")
    }
    dPort = buf[len(buf)-2:]
    return fmt.Sprintf("%s:%d", net.IP(dIP).String(), binary.BigEndian.Uint16(dPort)), nil
}
