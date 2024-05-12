package embed

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

func ExtractString(data []byte, off uint32) string {
	size := binary.BigEndian.Uint16(data[off : off+2])
	return string(data[int(off)+2 : int(off)+2+int(size)])
}

func ExtractStructs[T any](data []byte) []*T {
	var def T
	size := int(unsafe.Sizeof(def))
	length := len(data) / size
	out := make([]*T, length)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	for i := 0; i < length; i++ {
		meta := (*T)(unsafe.Pointer(hdr.Data + uintptr(i*size)))
		out[i] = meta
	}
	return out
}

func EmbedStructs[T any](data []*T) []byte {
	buf := bytes.NewBuffer(nil)

	var def T
	size := int(unsafe.Sizeof(def))
	slice := []byte{}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sh.Cap = size
	sh.Len = size

	for _, n := range data {
		sh.Data = uintptr(unsafe.Pointer(n))
		buf.Write(slice)
	}

	return buf.Bytes()
}
