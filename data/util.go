package data

import (
	"encoding/binary"
	"io"
)

const hextable = "0123456789ABCDEF"

//faster than fmt and need upper case!
func b2h(h []byte) []byte {
	b := make([]byte, len(h)*2)
	for i, v := range h {
		b[i*2] = hextable[v>>4]
		b[i*2+1] = hextable[v&0x0f]
	}
	return b
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func abs(a int64) uint64 {
	if a < 0 {
		return uint64(-a)
	}
	return uint64(a)
}

func write(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func read(r Reader, dest interface{}) error {
	return binary.Read(r, binary.BigEndian, dest)
}
