package read

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
		r := bytes.NewReader(b)
		v, _, err := VarInt(r)
		if err != nil {
			t.Error(err)
		}
		if v != value {
			t.Errorf("expected %d got %d, bytes: %v", value, v, b)
		}
	}
}
