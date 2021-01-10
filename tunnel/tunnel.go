package tunnel

import (
	"bufio"
	"fmt"
	"net"
)

type CiphAction func([]byte)

func Transport(dst net.Conn, src net.Conn, action CiphAction) {
	defer dst.Close()
	defer src.Close()
	reader := bufio.NewReader(src)
	buf := make([]byte, 4096)
	for {
		nsize, err := reader.Read(buf)
		if err != nil {
			return
		}
		if nsize > 0 {
			action(buf[:nsize])
			start := 0
			for {
				size, err := dst.Write(buf[start:nsize])
				if err != nil {
					fmt.Println("write error", err)
					return
				}
				start += size
				if start == nsize {
					break
				}
			}
		}
	}
}
