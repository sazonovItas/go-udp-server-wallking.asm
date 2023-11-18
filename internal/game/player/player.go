package player

import (
	"fmt"
	"net"
	"time"

	convertBytes "asm-game/server/internal/convertbytes"
	game "asm-game/server/internal/game"
)

type Info struct {
	// Name
	Name string

	// Player position
	Pos    game.Vec3
	Angles game.Vec3
	Size   game.Vec3
	Tex    int32
}

func (pli *Info) String() string {
	return fmt.Sprintf(
		"Name: %s\nPlayer position - %s\nAngles - %s\nSize - %s\nTex - %d",
		pli.Name,
		pli.Pos.String(),
		pli.Angles.String(),
		pli.Size.String(),
		pli.Tex,
	)
}

func (pli *Info) Update(data []byte) {
	pli.Pos.ConvertFromBytes(data[:12])
	pli.Angles.ConvertFromBytes(data[12:24])
	pli.Size.ConvertFromBytes(data[24:36])
	pli.Tex, _ = convertBytes.ByteSliceToT[int32](data[36:40])
}

func (pli *Info) ConvertToBytes() []byte {
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
	buf = append(buf, convertBytes.TToByteSlice[int32](pli.Tex)...)
	return buf
}

func New(addr string, data []byte) *Player {
	listenAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		fmt.Printf("Error occured while resolveing %s: %s", addr, err)
		return nil
	}

	buf := make([]byte, 0, 256)
	for _, v := range data {
		if v != 0 {
			buf = append(buf, v)
		}
	}

	return &Player{
		Addr:   listenAddr,
		Uptime: time.Now(),
		Info:   &Info{Name: string(buf)},
	}
}

type Player struct {
	// address and uptime
	Addr          *net.UDPAddr
	Uptime        time.Time
	SessionUpTime int32

	// Player data
	Info *Info
}

func (pl *Player) String() string {
	return fmt.Sprintf(
		"Addr: %s last update: %s\n%s",
		pl.Addr.String(),
		pl.Uptime.Format(time.UnixDate),
		pl.Info.String(),
	)
}
