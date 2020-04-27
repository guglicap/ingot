// package read implements methods to decode minecraft protocol data types
// from a io.Reader
package read

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func VarInt(r io.Reader) (v int32, n int, err error) {
	numRead := 0
	result := int32(0)
	var read byte
	for {
		read, err = UByte(r)
		if err != nil {
			return 0, numRead, err
		}
		v := int32(read & 0b01111111)
		result |= v << (7 * numRead)
		numRead++
		if numRead > 5 {
			return 0, numRead, errors.New("varint too big")
		}
		if (read & 0b10000000) == 0 {
			break
		}
	}
	return result, numRead, err
}

func VarLong(r io.Reader) (v int64, n int, err error) {
	numRead := 0
	result := int64(0)
	var read byte
	for {
		read, err = UByte(r)
		if err != nil {
			return 0, numRead, err
		}
		v := int64(read & 0b01111111)
		result |= v << (7 * numRead)
		numRead++
		if numRead > 10 {
			return 0, numRead, errors.New("varlong too big")
		}
		if (read & 0b10000000) == 0 {
			break
		}
	}
	return result, numRead, err
}

func Bool(r io.Reader) (bool, error) {
	b, err := UByte(r)
	if err != nil {
		return false, err
	}
	if b == 0x01 {
		return true, nil
	}
	if b == 0x00 {
		return false, nil
	}
	return false, errors.New("invalid bool encoding")
}

func UByte(r io.Reader) (uint8, error) {
	x := []byte{0x00}
	_, err := io.ReadFull(r, x)
	return x[0], err
}

func Byte(r io.Reader) (int8, error) {
	b, err := UByte(r)
	return int8(b), err
}

func UShort(r io.Reader) (uint16, error) {
	var s uint16
	err := binary.Read(r, binary.BigEndian, &s)
	return s, err
}

func Short(r io.Reader) (int16, error) {
	s, err := UShort(r)
	return int16(s), err
}

func Int(r io.Reader) (int32, error) {
	var i int32
	err := binary.Read(r, binary.BigEndian, &i)
	return i, err
}

func Long(r io.Reader) (int64, error) {
	var i int64
	err := binary.Read(r, binary.BigEndian, &i)
	return i, err
}

func Float(r io.Reader) (float32, error) {
	var f float32
	err := binary.Read(r, binary.BigEndian, &f)
	return f, err
}

func Double(r io.Reader) (float64, error) {
	var f float64
	err := binary.Read(r, binary.BigEndian, &f)
	return f, err
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
		return "", fmt.Errorf("invalid string length: wanted %d, got %d", length, n)
	}
	return string(buf), nil
}
