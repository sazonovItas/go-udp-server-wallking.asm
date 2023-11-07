package game

import (
	"fmt"

	convertBytes "asm-game/server/internal/convertbytes"
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func (v3 *Vec3) String() string {
	return fmt.Sprintf("x: %7.2f y: %7.2f z: %7.2f", v3.X, v3.Y, v3.Z)
}

func (v3 *Vec3) ConvertFromBytes(data []byte) {
	buf := convertBytes.ByteSliceToTSlice[float32](data)
	v3.X, v3.Y, v3.Z = buf[0], buf[1], buf[2]
}

type Vec4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (v4 *Vec4) String() string {
	return fmt.Sprintf("x: %7.2f y: %7.2f z: %7.2f w: %7.2f", v4.X, v4.Y, v4.Z, v4.W)
}

func (v4 *Vec4) ConvertFromBytes(data []byte) {
	buf := convertBytes.ByteSliceToTSlice[float32](data)
	v4.X, v4.Y, v4.Z, v4.W = buf[0], buf[1], buf[2], buf[3]
}
