package main


import "ladder/local"


func main() {
    loc := local.NewLocal(":1919", "your-remote-address:2020")
    loc.Serve()
}
