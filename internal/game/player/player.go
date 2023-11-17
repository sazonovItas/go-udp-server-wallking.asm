package player

import (
	"fmt"
	"net"
	"time"

	convertBytes "asm-game/server/internal/convertbytes"
	game "asm-game/server/internal/game"
)

type PlayerInfo struct {
	// Name
	Name string

	// Player position
	Pos    game.Vec3
	Angles game.Vec3
	Size   game.Vec3
}

func (pli *PlayerInfo) String() string {
	return fmt.Sprintf(
		"Name: %s\nPlayer position - %s\nAngles - %s\nSize - %s",
		pli.Name,
		pli.Pos.String(),
		pli.Angles.String(),
		pli.Size.String(),
	)
}

func (pli *PlayerInfo) Update(data []byte) {
	pli.Pos.ConvertFromBytes(data[:12])
	pli.Angles.ConvertFromBytes(data[12:24])
	pli.Size.ConvertFromBytes(data[24:36])
}

func (pli *PlayerInfo) ConvertToBytes() []byte {
	var buf []byte
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.X)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.Y)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.Z)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Angles.X)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Angles.Y)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Angles.Z)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Size.X)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Size.Y)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Size.Z)...)
	return buf
}

func New(addr string, data []byte) *Player {
	listenAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		fmt.Printf("Error occured while resolveing %s: %s", addr, err)
		return nil
	}

	return &Player{
		Addr:   listenAddr,
		Uptime: time.Now(),
		Info:   &PlayerInfo{Name: string(data)},
	}
}

type Player struct {
	// address and uptime
	Addr          *net.UDPAddr
	Uptime        time.Time
	SessionUpTime int32

	// Player data
	Info *PlayerInfo
}

func (pl *Player) String() string {
	return fmt.Sprintf(
		"Addr: %s last update: %s\n%s",
		pl.Addr.String(),
		pl.Uptime.Format(time.UnixDate),
		pl.Info.String(),
	)
}
