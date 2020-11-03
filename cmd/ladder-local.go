package main

import (
    "gosocks/local"
)


func main() {
    loc := local.NewLocal(":1919" ,":2020")
    loc.Serve()
}
