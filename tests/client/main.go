package main

import (
	"log"
	"net"
	"os"

	convertBytes "asm-game/server/internal/convertbytes"
)

func main() {
	port := "60105"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	serverAddr, err := net.ResolveUDPAddr("udp4", "192.168.1.255:8829")
	if err != nil {
		log.Fatalf("Resolving udp addr occured: %s", err)
	}

	listenAddr, err := net.ResolveUDPAddr("udp4", ":"+port)
	if err != nil {
		log.Fatalf("Resolving udp addr occured: %s", err)
	}

	listen, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		log.Fatalf("Listen udp addr occured: %s", err)
	}
	defer listen.Close()

	// join test
	buf := make([]byte, 0, 256)
	buf = append(buf, convertBytes.TToByteSlice[int32](256)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, convertBytes.TToByteSlice[int32](1)...)
	buf = append(buf, []byte("Alex")...)
	buf = buf[:256]
	n, err := listen.WriteToUDP(buf, serverAddr)
	if err != nil {
		log.Fatalf("Error %s", err.Error())
	}
	log.Printf("Send data size: %d bytes", n)

	n, addr, err := listen.ReadFromUDP(buf)
	if err != nil {
		log.Fatalf("Error %s", err.Error())
	}
	log.Printf("Msg from %s: %s", addr.String(), buf[:n])
	uptime, ok := convertBytes.ByteSliceToT[int32](buf[14:])
	if !ok && string(buf[12:14]) == "Ok" {
		return
	}

	buf = make([]byte, 0, 256)
	buf = append(buf, convertBytes.TToByteSlice[int32](256)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, []byte{00, 00, 00, 00}...)
	buf = append(buf, convertBytes.TToByteSlice[int32](uptime)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](360.22)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](180.43)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](23.01)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](12.92)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](23.01)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](12.92)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](12.92)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](23.01)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](12.92)...)
	buf = buf[:256]
	_, err = listen.WriteToUDP(
		buf,
		serverAddr,
	)
	if err != nil {
		log.Fatalf("Write to udp socket occured: %s", err)
	}

	buf = make([]byte, 256)
	n, addr, err = listen.ReadFromUDP(buf)
	if err != nil {
		log.Fatalf("Cannot read from udp addr occured: %s", err)
	}

	uptime, _ = convertBytes.ByteSliceToT[int32](buf[12:16])
	log.Printf("Uptime: %d", uptime)
	cntPlayers, _ := convertBytes.ByteSliceToT[int32](buf[16:20])
	log.Printf("Count players: %d", cntPlayers)
	flbuf := convertBytes.ByteSliceToTSlice[float32](buf[20:n])
	for _, v := range flbuf {
		log.Printf("%7.2f ", v)
	}
	log.Printf("Udp message from %s: %s\n", addr, buf[:n])
}
