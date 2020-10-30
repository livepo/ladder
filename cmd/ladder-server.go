package main
import (
    "ladder/server"
)


func main() {
    srv := server.NewServer(":2020")
    srv.Serve()
}
