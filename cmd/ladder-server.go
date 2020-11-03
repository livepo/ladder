package main

import (
    "gosocks/server"
)


func main() {
    serv := server.NewServer(":2020")
    serv.Serve()
}
