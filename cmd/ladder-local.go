package main

import (
	"flag"
	"fmt"
	"gosocks/cipher"
	"gosocks/local"
)

func main() {
	var listenaddr, remoteaddr string
	var secret string
	flag.StringVar(&listenaddr, "listen", ":1919", "local listen address")
	flag.StringVar(&remoteaddr, "remote", ":2020", "remote connect address")
	flag.StringVar(&secret, "secret", cipher.LoadSecret(), "secret recorded under server's start directory named .gosocks")
	flag.Parse()
	if len(secret) == 0 || len(listenaddr) == 0 || len(remoteaddr) == 0 {
		fmt.Println("usage: ./ladder-local --listenaddr :1919 --remoteaddr :2020 --secret [secret_token]")
		return
	}
	cipher.UpdateSecret(secret)
	loc := local.NewLocal(listenaddr, remoteaddr)
	fmt.Printf("local  address: %s\n", listenaddr)
	fmt.Printf("remote address: %s\n", remoteaddr)
	fmt.Printf("secret: %s\n", secret)
	loc.Serve()
}
