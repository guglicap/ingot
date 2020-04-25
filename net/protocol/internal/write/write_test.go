package write

import (
	"bytes"
	"testing"
)

var varIntTest = map[int32][]byte{
	0:           {0x00},
	1:           {0x01},
	2:           {0x02},
	127:         {0x7f},
	128:         {0x80, 0x01},
	255:         {0xff, 0x01},
	2147483647:  {0xff, 0xff, 0xff, 0xff, 0x07},
	-1:          {0xff, 0xff, 0xff, 0xff, 0x0f},
	-2147483648: {0x80, 0x80, 0x80, 0x80, 0x08},
}

func TestVarInt(t *testing.T) {
	for value, b := range varIntTest {
		buf := bytes.NewBuffer([]byte{})
		_, err := VarInt(value, buf)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(b, buf.Bytes()) {
			t.Errorf("wrong encoding, wanted: %v have %v", b, buf.Bytes())
		}
	}
}
