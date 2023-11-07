package player

import (
	"fmt"
	"net"
	"time"

	game "asm-game/server/game"
	convertBytes "asm-game/server/internal/convertbytes"
)

type PlayerInfo struct {
	// Name
	Name string

	// Player position
	Pos game.Vec3
	Yaw float32
}

func (pli *PlayerInfo) String() string {
	return fmt.Sprintf(
		"Name: %s\nPlayer position - %s\nYaw angle - %7.2f",
		pli.Name,
		pli.Pos.String(),
		pli.Yaw,
	)
}

func (pli *PlayerInfo) Update(data []byte) {
	pli.Pos.ConvertFromBytes(data[:12])
	v, ok := convertBytes.ByteSliceToT[float32](data[12:16])
	if !ok {
		fmt.Printf("Problem with converting float32 from byte slice\n")
		return
	}
	pli.Yaw = v
}

func (pli *PlayerInfo) ConvertToBytes() []byte {
	var buf []byte
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.X)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.Y)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Pos.Z)...)
	buf = append(buf, convertBytes.TToByteSlice[float32](pli.Yaw)...)
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
	Addr   *net.UDPAddr
	Uptime time.Time

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
