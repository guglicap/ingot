// package write implements methods to encode encode minecraft protocol data types
// to a io.Writer
package write

import (
	"encoding/binary"
	"io"
)

func VarInt(value int32, w io.Writer) (int, error) {
	i := 0
	v := uint32(value)
	for {
		b := v & 0x7F
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		err := UByte(uint8(b), w)
		if err != nil {
			return i, err
		}
		i++
		if v == 0 {
			break
		}
	}
	return i, nil
}

func VarLong(value int64, w io.Writer) (int, error) {
	i := 0
	for {
		temp := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			temp |= 0b10000000
		}
		err := UByte(temp, w)
		if err != nil {
			return i, err
		}
		i++
		if value == 0 {
			break
		}
	}
	return i, nil
}

func Bool(b bool, w io.Writer) error {
	a := byte(0)
	if b {
		a = 1
	}
	return UByte(a, w)
}

func UByte(b byte, w io.Writer) error {
	_, err := w.Write([]byte{b})
	return err
}

func Byte(b int8, w io.Writer) error {
	return UByte(uint8(b), w)
}

func UShort(v uint16, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}

func Short(v int16, w io.Writer) error {
	return UShort(uint16(v), w)
}

func Int(v int32, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}

func Long(v int64, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}

func Float(v float32, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}

func Double(v float64, w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}

func String(v string, w io.Writer) error {
	_, err := VarInt(int32(len(v)), w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(v))
	return err
}
