package server

import (
	"asm-game/server/internal/game/player"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/exp/maps"

	convertBytes "asm-game/server/internal/convertbytes"
)

// data [0..255]:
// FOR ALL:
// [0..3] - check sum;
// [4..11] - WallKing
// [12..15] - state of the player;
// JOIN:
// [16..255] - name
// UPDATE:
// [16..19] - session upTime
// [20..55] - position, angles and size of player;
// EXIT:
const (
	PlJoin   int32 = 1
	PlUpdate int32 = 0
	PlExit   int32 = -1
)

func New(log *slog.Logger, address string, port string) *Server {
	listenAddr, err := net.ResolveUDPAddr("udp4", address+":"+port)
	if err != nil {
		fmt.Printf("Error occured while resolve udp addr: %s", err)
		os.Exit(1)
	}

	return &Server{
		Addr:       listenAddr,
		ListenCon:  nil,
		Logger:     log,
		UpSendTime: time.Now(),
		session:    &Session{CntPl: 0, SessionTime: 0, Players: map[string]*player.Player{}},
	}
}

type Server struct {
	// Server's address and listen connection
	Addr      *net.UDPAddr
	ListenCon *net.UDPConn
	Logger    *slog.Logger

	UpSendTime time.Time
	session    *Session
}

func (sv *Server) Up() {
	const op = "server.server.go"

	listen, err := net.ListenUDP("udp4", sv.Addr)
	if err != nil {
		sv.Logger.Error("Error cannot listen udp", op, err.Error())
		os.Exit(1)
	}

	sv.Logger.Info(
		"Server started",
		slog.String("address", sv.Addr.String()),
		slog.String("time", time.Now().Format(time.UnixDate)),
	)
	sv.ListenCon = listen
}

func (sv *Server) Down() {
	sv.Logger.Info("server is down", slog.String("time", time.Now().Format(time.UnixDate)))
	err := sv.ListenCon.Close()
	if err != nil {
		sv.Logger.Error("close listen connection", slog.String("error", err.Error()))
		return
	}
}

func (sv *Server) UpdatePlConn(addr string, data []byte) {
	// check length of message
	if len(data) < 4 {
		sv.Logger.Debug("Wrong format of data")
		return
	}

	// check checkSum
	checkSum, ok := convertBytes.ByteSliceToT[int32](data[:4])
	if !ok || checkSum != (int32)(len(data)) {
		sv.Logger.Debug(
			"Wrong format of data",
			slog.String("address", addr),
			slog.String("checksum", "wrong checksum"),
			slog.Int("checksum", (int)(checkSum)),
			slog.Int("len of data", len(data)),
		)
		return
	}

	// check game name
	if len(data) < 11 && string(data[4:11]) != "WallKing" {
		sv.Logger.Debug(
			"Wrong format of data",
			slog.String("WallKing", "wrong sync game name"),
		)
		return
	}

	// check state of the player
	state, ok := convertBytes.ByteSliceToT[int32](data[12:16])
	if !ok || (state != -1) && (state != 1) && (state != 0) {
		sv.Logger.Debug(
			"Wrong format of data",
			slog.String("state", "wrong state"),
		)
		return
	}

	switch state {
	case PlJoin:
		if len(data) < 17 {
			sv.Logger.Debug("Wrong join message", slog.String("msg", string(data)))
			return
		}
		sv.session.NewPlayer(addr, data[16:])
		go sv.SendToNew(addr)
	case PlUpdate:
		sv.session.Lock()
		sv.session.SessionTime++
		sv.session.Unlock()
		if len(data) < 56 {
			sv.Logger.Debug("Wrong update message", slog.String("msg", string(data)))
			return
		}
		sv.session.UpdatePlayer(addr, data[16:])
		go sv.SendToAll()
	case PlExit:
		sv.session.ExitPlayer(addr)
	}
}

func (sv *Server) SendToNew(addr string) {
	pl, ok := sv.session.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		return
	}

	sv.Logger.Info("Server sending join accepts", slog.String("address", addr))
	buf := make([]byte, 0, 256)
	buf = append(buf, convertBytes.TToByteSlice[int32](256)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, []byte("Ok")...)
	buf = append(buf, convertBytes.TToByteSlice[int32](sv.session.SessionTime)...)
	buf = buf[:256]
	_, err := sv.ListenCon.WriteToUDP(buf, pl.Addr)
	if err != nil {
		sv.Logger.Error("Do not send join accepts to player", slog.String("address", addr))
		return
	}
}

func (sv *Server) SendToAll() {
	sv.session.Lock()
	players := maps.Values(sv.session.Players)
	sv.session.Unlock()

	for _, pl := range players {
		go sv.sendToPlayer(pl, players)
	}
}

func (sv *Server) sendToPlayer(pl *player.Player, players []*player.Player) {
	buf := make([]byte, 0, 256)

	buf = append(buf, convertBytes.TToByteSlice[int32](256)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, convertBytes.TToByteSlice[int32](sv.session.SessionTime)...)
	buf = append(buf, convertBytes.TToByteSlice[int32]((int32)(len(players))-1)...)
	for _, v := range players {
		if v != pl {
			buf = append(buf, v.Info.ConvertToBytes()...)
		}
	}

	buf = buf[:256]
	_, err := sv.ListenCon.WriteToUDP(buf, pl.Addr)
	if err != nil {
		sv.Logger.Error("Error to send data to Player", slog.String("address", pl.Addr.String()))
		return
	}
}
