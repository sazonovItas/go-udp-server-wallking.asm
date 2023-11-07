package convertBytes

import (
	"reflect"
	"unsafe"
)

type Number interface {
	int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64 | float32 | float64
}

func ByteSliceToT[T Number](src []byte) (T, bool) {
	var value T
	if len(src) == 0 {
		return value, false
	}

	ptr := unsafe.Pointer(&src[0])
	return *((*T)(ptr)), true
}

func ByteSliceToTSlice[T Number](src []byte) []T {
	if len(src) == 0 {
		return nil
	}

	var value T
	length := len(src) / (int)(unsafe.Sizeof(reflect.TypeOf(value)))
	ptr := unsafe.Pointer(&src[0])
	// It is important to keep in mind that Go garbage collector
	// will not interact with this data, and that if src freed,
	// the behavior of any Go code using the slice is nondeterministic
	return (*[1 << 24]T)((*[1 << 24]T)(ptr))[:length:length]
}

func TToByteSlice[T Number](src T) []byte {
	length := (int)(unsafe.Sizeof(reflect.TypeOf(src)))

	buf := make([]byte, length)
	ptr := unsafe.Pointer(&buf[0])
	*((*T)(ptr)) = src

	return buf
}

func TSliceToByteSlice[E ~[]T, T Number](src E) []byte {
	if len(src) == 0 {
		return nil
	}
	length := (int)(unsafe.Sizeof(reflect.TypeOf(src[0])))

	buf := make([]byte, length*len(src))
	for i := 0; i < len(src)*length; i += length {
		ptr := unsafe.Pointer(&buf[i])
		*((*T)(ptr)) = src[i/length]
	}

	return buf
}
