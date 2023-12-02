package server

import (
	"asm-game/server/internal/config"
	"asm-game/server/internal/game/player"
	"asm-game/server/internal/server/msgqueue"
	"asm-game/server/internal/server/session"
	"log/slog"
	"net"
	"os"
	"time"

	"golang.org/x/exp/maps"

	convertBytes "asm-game/server/internal/convertbytes"
)

// data [0..255]:
// FOR ALL:
// [0..1] - check sum;
// [2..9] - WallKing
// [10..10] - state of the player;
// JOIN:
// [11..511] - name
// UPDATE:
// [11..14] - session upTime
// [15..50] - position, angles and size of player;
// [51..54] - ambient texture
// [55..68] - diffuse texture
// [69..62] - specular texture
// [63..66] - shininess
// EXIT:
const (
	PlJoin   int8 = 1
	PlUpdate int8 = 0
	PlExit   int8 = -1

	// MSGSize - msg size
	MSGSize int16 = 512
)

func New(log *slog.Logger, cfgServer config.UDPServer) *Server {
	listenAddr, err := net.ResolveUDPAddr("udp4", cfgServer.Address+":"+cfgServer.Port)
	if err != nil {
		log.Error("Error occurred while resolve udp addr: %s", err.Error())
		os.Exit(1)
	}

	return &Server{
		msgQueue:    msgqueue.New(),
		Addr:        listenAddr,
		ListenCon:   nil,
		Logger:      log,
		SendTimeout: cfgServer.Timeout,
		PlTimeout:   cfgServer.IdleTimeout,
		UpSendTime:  time.Now(),
		session:     &session.Session{CntPl: 0, SessionTime: 0, Players: map[string]*player.Player{}},
	}
}

type Server struct {
	// Message queue
	msgQueue *msgqueue.MsgQueue

	// Server's address and listen connection
	Addr      *net.UDPAddr
	ListenCon *net.UDPConn
	Logger    *slog.Logger

	SendTimeout time.Duration
	PlTimeout   time.Duration

	UpSendTime time.Time
	session    *session.Session
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

func (sv *Server) HandleMsg(addr string, data []byte) {
	// check length of message
	if len(data) < 2 {
		sv.Logger.Debug("Wrong format of data")
		return
	}

	// check checkSum
	checkSum, ok := convertBytes.ByteSliceToT[int16](data[:4])
	if !ok || checkSum != (int16)(len(data)) {
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
	if len(data) < 10 && string(data[2:10]) != "WallKing" {
		sv.Logger.Debug(
			"Wrong format of data",
			slog.String("WallKing", "wrong sync game name"),
		)
		return
	}

	// check state of the player
	state, ok := convertBytes.ByteSliceToT[int8](data[10:11])
	if !ok || (state != PlExit) && (state != PlJoin) && (state != PlUpdate) {
		sv.Logger.Debug(
			"Wrong format of data",
			slog.String("state", "wrong state"),
		)
		return
	}

	switch state {
	case PlJoin:
		if len(data) < 11 {
			sv.Logger.Debug("Wrong join message", slog.String("msg", string(data)))
			return
		}
		sv.session.NewPlayer(addr, data[11:])
		sv.SendToNew(addr)
	case PlUpdate:
		if len(data) < 66 {
			sv.Logger.Debug("Wrong update message", slog.String("msg", string(data)))
			return
		}
		sv.session.UpdatePlayer(addr, data[11:])
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
	buf := make([]byte, 0, MSGSize)
	buf = append(buf, convertBytes.TToByteSlice[int16](MSGSize)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, convertBytes.TToByteSlice[int16](200)...)
	buf = append(buf, convertBytes.TToByteSlice[int32](sv.session.SessionTime)...)
	buf = buf[:MSGSize]
	msg := msgqueue.Message{
		Addr: pl.Addr,
		Msg:  buf,
	}

	sv.msgQueue.Lock()
	sv.msgQueue.Queue.Enqueue(msg)
	sv.msgQueue.Unlock()
}

func (sv *Server) SendToAll() {
	sv.session.Lock()
	players := maps.Values(sv.session.Players)
	sv.session.Unlock()

	sv.msgQueue.Lock()
	for _, pl := range players {
		if sv.UpSendTime.Sub(pl.Uptime).Seconds() > sv.PlTimeout.Seconds() {
			sv.session.ExitPlayer(pl.Addr.String())
		}
		if pl.Updated {
			sv.sendToPlayer(pl, players)
		}
	}
	sv.msgQueue.Unlock()
}

func (sv *Server) sendToPlayer(pl *player.Player, players []*player.Player) {
	buf := make([]byte, 0, MSGSize)

	buf = append(buf, convertBytes.TToByteSlice[int16](MSGSize)...)
	buf = append(buf, []byte("WallKing")...)
	buf = append(buf, convertBytes.TToByteSlice[int32](sv.session.SessionTime)...)
	buf = append(buf, convertBytes.TToByteSlice[int32]((int32)(len(players))-1)...)
	for _, v := range players {
		if v.Addr.String() != pl.Addr.String() && v.Updated {
			buf = append(buf, v.Info.ConvertToBytes()...)
		}
	}

	buf = buf[:MSGSize]
	msg := msgqueue.Message{
		Addr: pl.Addr,
		Msg:  buf,
	}
	sv.msgQueue.Queue.Enqueue(msg)
}

func (sv *Server) SendMessages() {

	sv.UpSendTime = time.Now()

	for {
		if time.Now().Sub(sv.UpSendTime).Milliseconds() > sv.SendTimeout.Milliseconds() {
			sv.UpSendTime = time.Now()
			sv.session.Lock()
			sv.session.SessionTime++
			sv.session.Unlock()
			go sv.SendToAll()
		}
		sv.msgQueue.Lock()
		for !sv.msgQueue.Queue.Empty() {
			msg := sv.msgQueue.Queue.Dequeue()
			_, err := sv.ListenCon.WriteToUDP(msg.Msg, msg.Addr)
			if err != nil {
				sv.Logger.Error(
					"Error to send data to Player",
					slog.String("address", msg.Addr.String()),
				)
				sv.msgQueue.Unlock()
				return
			}
		}
		sv.msgQueue.Unlock()
	}
}
