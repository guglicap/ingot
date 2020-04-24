package read

import (
	"errors"
	"io"
)

func VarInt(r io.Reader) (v int64, n int, err error) {
	result := int64(0)
	for {
		var x byte
		x, err = Byte(r)
		if err != nil {
			return
		}
		value := x & 0x7f
		result |= int64(uint(value) << uint(7*n))
		n++
		if n > 5 {
			err = errors.New("varint too big")
			return
		}

		if (x & 0x80) == 0 {
			break
		}
	}
	v = result
	return
}

func Byte(r io.Reader) (byte, error) {
	x := []byte{0x00}
	_, err := io.ReadFull(r, x)
	return x[0], err
}

func String(r io.Reader) (string, error) {
	length, _, err := VarInt(r)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	if n != int(length) {
		// TODO: logging
	}
	return string(buf), nil
}
