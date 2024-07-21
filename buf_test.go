package buf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Buf(t *testing.T) {
	items := [][]struct {
		ret    []byte
		isVar  bool
		str    string
		bytes  []byte
		num8   int8
		num16  int16
		num32  int32
		num64  int64
		numVar int64
	}{
		{{numVar: 0xfffffffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0f}}},
		{{numVar: 0xffffffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}}},
		{{numVar: 0xfffffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x07}}},
		{{numVar: 0xffffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f}}},
		{{numVar: 0xfffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x03}}},
		{{numVar: 0xffffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}}},
		{{numVar: 0xfffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x01}}},
		{{numVar: 0xffffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0xff, 0x0f}}},
		{{numVar: 0xfffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0x7f}}},
		{{numVar: 0xffffff, isVar: true, ret: []byte{0xff, 0xff, 0xff, 0x07}}},
		{{numVar: 0xfffff, isVar: true, ret: []byte{0xff, 0xff, 0x3f}}},
		{{numVar: 0xffff, isVar: true, ret: []byte{0xff, 0xff, 0x03}}},
		{{numVar: 0xfff, isVar: true, ret: []byte{0xff, 0x1f}}},
		{{numVar: 0xff, isVar: true, ret: []byte{0xff, 0x1}}},
		{{numVar: 0xf, isVar: true, ret: []byte{0xf}}},
		{{numVar: 0x0, isVar: true, ret: []byte{0x0}}},
		{{numVar: 0x80, isVar: true, ret: []byte{0x80, 0x01}}},

		{{num64: 0xfffffffffffffff, isVar: false, ret: []byte{0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xffffffffffffff, isVar: false, ret: []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xfffffffffffff, isVar: false, ret: []byte{0x00, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xffffffffffff, isVar: false, ret: []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xfffffffffff, isVar: false, ret: []byte{0x00, 0x00, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xffffffffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xfffffffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x0f, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xffffffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}}},
		{{num64: 0xfffffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x0f, 0xff, 0xff, 0xff}}},
		{{num64: 0xffffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff}}},
		{{num64: 0xfffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0xff, 0xff}}},
		{{num64: 0xffff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff}}},
		{{num64: 0xfff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0xff}}},
		{{num64: 0xff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff}}},
		{{num64: 0xf, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f}}},

		{{num32: 0xfffffff, isVar: false, ret: []byte{0x0f, 0xff, 0xff, 0xff}}},
		{{num32: 0xffffff, isVar: false, ret: []byte{0x00, 0xff, 0xff, 0xff}}},
		{{num32: 0xfffff, isVar: false, ret: []byte{0x00, 0x0f, 0xff, 0xff}}},
		{{num32: 0xffff, isVar: false, ret: []byte{0x00, 0x00, 0xff, 0xff}}},
		{{num32: 0xfff, isVar: false, ret: []byte{0x00, 0x00, 0x0f, 0xff}}},
		{{num32: 0xff, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0xff}}},
		{{num32: 0xf, isVar: false, ret: []byte{0x00, 0x00, 0x00, 0x0f}}},

		{{num16: 0xfff, isVar: false, ret: []byte{0x0f, 0xff}}},
		{{num16: 0xff, isVar: false, ret: []byte{0x00, 0xff}}},
		{{num16: 0xf, isVar: false, ret: []byte{0x00, 0x0f}}},

		{
			{numVar: 0, isVar: true, ret: []byte{0}},
			{numVar: 127, isVar: true, ret: []byte{0, 0x7f}},
			{numVar: 128, isVar: true, ret: []byte{0, 0x7f, 0x80, 0x01}},
			{str: "abc", isVar: true, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c'}},
			{str: "123", isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3'}},
			{bytes: []byte{1, 2, 3}, isVar: true, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3}},
			{bytes: []byte{4, 5, 6}, isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3, 4, 5, 6}},
			{num32: 1, isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3, 4, 5, 6, 0x00, 0x00, 0x00, 0x01}},
			{num64: 1, isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3, 4, 5, 6, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
			{num16: 1, isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3, 4, 5, 6, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01}},
			{num8: 1, isVar: false, ret: []byte{0, 0x7f, 0x80, 0x01, 0x03, 'a', 'b', 'c', '1', '2', '3', 3, 1, 2, 3, 4, 5, 6, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01}},
		},
	}

	for _, section := range items {
		buf := new(bytes.Buffer)
		pbW := NewWriter(buf)
		for _, item := range section {
			if item.isVar {
				if len(item.bytes) > 0 {
					pbW.WriteVarBytes(item.bytes)
				} else if len(item.str) > 0 {
					pbW.WriteVarString(item.str)
				} else if item.numVar >= 0 {
					pbW.WriteVarInt(uint64(item.numVar))
				}
			} else {
				if len(item.bytes) > 0 {
					pbW.WriteBytes(item.bytes)
				} else if len(item.str) > 0 {
					pbW.WriteString(item.str)
				} else if item.num8 > 0 {
					pbW.WriteInt8(uint8(item.num8))
				} else if item.num16 > 0 {
					pbW.WriteInt16(uint16(item.num16))
				} else if item.num32 > 0 {
					pbW.WriteInt32(uint32(item.num32))
				} else if item.num64 > 0 {
					pbW.WriteInt64(uint64(item.num64))
				}
			}
			bufByte := buf.Bytes()
			assert.Equal(t, item.ret, bufByte)
		}
		pbR := NewReader(bytes.NewBuffer(buf.Bytes()))
		for _, item := range section {
			if item.isVar {
				if len(item.bytes) > 0 {
					data, _, err := pbR.ReadVarBytes()
					assert.Nil(t, err)
					assert.Equal(t, item.bytes, data)
				} else if len(item.str) > 0 {
					data, _, err := pbR.ReadVarString()
					assert.Nil(t, err)
					assert.Equal(t, item.str, data)
				} else if item.numVar >= 0 {
					data, _, err := pbR.ReadVarInt()
					assert.Nil(t, err)
					assert.Equal(t, uint64(item.numVar), data)
				}
			} else {
				if len(item.bytes) > 0 {
					data, err := pbR.ReadBytes(uint64(len(item.bytes)))
					assert.Nil(t, err)
					assert.Equal(t, item.bytes, data)
				} else if len(item.str) > 0 {
					data, err := pbR.ReadString(uint64(len(item.str)))
					assert.Nil(t, err)
					assert.Equal(t, item.str, data)
				} else if item.num8 > 0 {
					data, err := pbR.ReadInt8()
					assert.Nil(t, err)
					assert.Equal(t, uint8(item.num8), data)
				} else if item.num16 > 0 {
					data, err := pbR.ReadInt16()
					assert.Nil(t, err)
					assert.Equal(t, uint16(item.num16), data)
				} else if item.num32 > 0 {
					data, err := pbR.ReadInt32()
					assert.Nil(t, err)
					assert.Equal(t, uint32(item.num32), data)
				} else if item.num64 > 0 {
					data, err := pbR.ReadInt64()
					assert.Nil(t, err)
					assert.Equal(t, uint64(item.num64), data)
				}
			}
		}
	}
}
