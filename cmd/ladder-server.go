package main

import (
    "flag"
    "fmt"
    "gosocks/cipher"
    "gosocks/server"
)


func main() {
    var listenaddr string
    flag.StringVar(&listenaddr, "listen", ":2020", "listen address")
    flag.Parse()
    serv := server.NewServer(listenaddr)
    fmt.Printf("server: %s\n", listenaddr)
    secret := cipher.InitPasswd()
    fmt.Printf("secret: %s\n", secret)
    serv.Serve()
}
