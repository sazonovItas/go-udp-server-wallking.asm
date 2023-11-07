package server

import (
	"asm-game/server/game/player"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	convertBytes "asm-game/server/internal/convertbytes"
)

// data [0..255] and [0..3] - state of the player
const (
	PL_JOIN   int32 = 1  // data [4..255]: [4..255] - name
	PL_UPDATE int32 = 0  // data [4..255]: [4..15] - position, [16..19] - yaw angle
	PL_EXIT   int32 = -1 // data [4..255]: whatever you want
)

func New() *Server {
	port, ip := os.Getenv("SERVER_PORT"), os.Getenv("SEVER_IP")
	listenAddr, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err != nil {
		fmt.Printf("Error occured while resolve udp addr: %s", err)
		os.Exit(1)
	}

	return &Server{
		Addr:      listenAddr,
		ListenCon: nil,
		session:   &Session{CntPl: 0, Players: map[string]*player.Player{}},
	}
}

type Server struct {
	// Server's address and listen connection
	Addr      *net.UDPAddr
	ListenCon *net.UDPConn

	session *Session
}

func (sv *Server) Up() {
	listen, err := net.ListenUDP("udp4", sv.Addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Server up at %s at %s\n",
		sv.Addr.IP.String()+":"+fmt.Sprintf("%d", sv.Addr.Port),
		time.Now().Format(time.UnixDate),
	)
	sv.ListenCon = listen
}

func (sv *Server) Down() {
	fmt.Printf("Server down at %s\n", time.Now().Format(time.UnixDate))
	sv.ListenCon.Close()
}

func (sv *Server) UpdatePlConn(addr string, data []byte) {
	if len(data) <= 4 {
		return
	}
	state, ok := convertBytes.ByteSliceToT[int32](
		data[:4],
	)
	log.Printf("State of player: %d\n", state)

	if !ok {
		log.Printf("converting is not possible\n")
	}

	switch state {
	case PL_JOIN:
		sv.session.NewPlayer(addr, data[4:])
		go sv.SendToNew(addr)
	case PL_UPDATE:
		sv.session.UpdatePlayer(addr, data[4:])
		go sv.SendToAll(addr)
	case PL_EXIT:
		sv.session.ExitPlayer(addr, data[4:])
	}
}

func (sv *Server) SendToNew(addr string) {
	pl, ok := sv.session.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		return
	}

	log.Printf("Server sending join acception to addr %s\n", addr)
	buf := make([]byte, 256)
	buf[0] = 0xFF
	_, err := sv.ListenCon.WriteToUDP(buf, pl.Addr)
	if err != nil {
		log.Printf("Do not send join acception to player %s\n", addr)
		return
	}
}

func (sv *Server) SendToAll(addr string) {
	pl, ok := sv.session.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		return
	}

	var buf []byte
	for k, v := range sv.session.Players {
		if k == addr {
			continue
		}

		buf = append(buf, v.Info.ConvertToBytes()...)
	}

	_, err := sv.ListenCon.WriteToUDP(buf, pl.Addr)
	if err != nil {
		log.Printf("Cannot send data to addr: %s\n", addr)
	}
}
