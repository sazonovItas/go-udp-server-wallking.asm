package main

import (
	"asm-game/server/cmd/server"
	"log"
)

func main() {
	udpServer := server.New()
	udpServer.Up()

	buf := make([]byte, 256)
	for {
		n, addr, err := udpServer.ListenCon.ReadFrom(buf)
		if err != nil {
			log.Printf("Reading from buffer occur: %s", err)
			continue
		}

		if string(buf[:n-1]) == "exit" {
			break
		}

		go udpServer.UpdatePlConn(addr.String(), buf[:n])
	}
	udpServer.Down()
}
