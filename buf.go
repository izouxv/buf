package buf

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Writer interface {
	WriteVarBytes(data []byte) error
	WriteVarString(str string) error
	WriteVarInt(num uint64) error

	WriteBytes(data []byte) error
	WriteString(str string) error

	WriteInt64(num uint64) error
	WriteInt32(num uint32) error
	WriteInt16(num uint16) error
	WriteInt8(num uint8) error

	CopyBytes(data io.Reader) error
}
type Reader interface {
	ReadVarBytes() (data []byte, varIntLen int, err error)
	ReadVarString() (str string, varIntLen int, err error)
	ReadVarInt() (num uint64, varIntLen int, err error)

	ReadBytes(length uint64) ([]byte, error)
	ReadString(length uint64) (string, error)

	ReadInt64() (uint64, error)
	ReadInt32() (uint32, error)
	ReadInt16() (uint16, error)
	ReadInt8() (uint8, error)
}
type WR interface {
	Writer
	Reader
}

func NewReaderWriter(out io.ReadWriter) WR {
	return newBuf(out, out)
}
func NewWriter(out io.Writer) Writer {
	return newBuf(nil, out)
}
func NewReader(in io.Reader) Reader {
	return newBuf(in, nil)
}
func newBuf(in io.Reader, out io.Writer) *pbBuf {
	return &pbBuf{
		in:     in,
		out:    out,
		data_9: make([]byte, 9),
	}
}

type pbBuf struct {
	in     io.Reader
	out    io.Writer
	data_9 []byte
}

func (s *pbBuf) WriteVarString(str string) error {
	return s.WriteVarBytes([]byte(str))
}
func (s *pbBuf) WriteVarBytes(data []byte) (err error) {
	err = s.WriteVarInt(uint64(len(data)))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
		return fmt.Errorf("data is empty")
	}
	return writeBytes(s.out, data)
}

func (s *pbBuf) WriteBytes(data []byte) error {
	return writeBytes(s.out, data)
}
func (s *pbBuf) WriteString(str string) error {
	return writeBytes(s.out, []byte(str))
}
func writeBytes(out io.Writer, data []byte) (err error) {
	cur := 0
	for {
		done, err := out.Write(data[cur:])
		if err != nil {
			return err
		}
		cur += done
		if cur == len(data) {
			break
		}
	}
	return
}

func (s *pbBuf) ReadVarString() (string, int, error) {
	data, size, err := s.ReadVarBytes()
	if err != nil {
		return "", 0, err
	}
	return string(data), size, nil
}
func (s *pbBuf) ReadVarBytes() (data []byte, varIntLen int, err error) {
	length, varIntLen, err := s.ReadVarInt()
	if err != nil {
		return nil, 0, err
	}
	if length == 0 {
		return
	}
	data = make([]byte, length)
	err = readBytes(s.in, data, length)
	return data, varIntLen, err
}
func (s *pbBuf) ReadString(length uint64) (string, error) {
	data, err := s.ReadBytes(length)
	return string(data), err
}
func (s *pbBuf) ReadBytes(length uint64) ([]byte, error) {
	data := make([]byte, length)
	err := readBytes(s.in, data, uint64(length))
	return data, err
}
func readBytes(in io.Reader, data []byte, length uint64) (err error) {
	var cur uint64 = 0
	for {
		done, err := in.Read(data[cur:])
		if err != nil {
			return err
		}
		if done < 0 {
			panic(-1)
		}
		cur += uint64(done)
		if cur == length {
			break
		}
	}
	return
}
func (s *pbBuf) CopyBytes(data io.Reader) error {
	return copyContentBuffer(s.out, data)
}
func (s *pbBuf) WriteVarInt(num2 uint64) (err error) {
	num := num2
	data := s.data_9[:0]
	offsize := 0
	for {
		tmp := byte(num & 0x7f)
		num = num >> 7

		if num > 0 {
			tmp = tmp | 0x80
		}
		offsize++
		data = append(data, tmp)
		if num == 0 {
			break
		}
		if num < 0 {
			panic(-1)
		}
	}

	err = writeBytes(s.out, data[:offsize])
	if err != nil {
		return
	}
	return nil
}

func (s *pbBuf) ReadVarInt() (num uint64, varIntLen int, err error) {
	data := s.data_9[:1]
	offset := 0
	for {
		_, err = s.in.Read(data)
		if err != nil {
			return
		}
		tmp := data[0]
		num = num + ((uint64(tmp) & 0x7f) << offset)
		offset += 7
		varIntLen++
		isEnd := tmp>>7 == 0
		if isEnd {
			break
		}
		if tmp>>7 < 0 {
			panic(-1)
		}
	}
	return
}

func (s *pbBuf) ReadInt64() (num uint64, err error) {
	err = readBytes(s.in, s.data_9[:8], 8)
	if err != nil {
		return
	}
	num = binary.BigEndian.Uint64(s.data_9)
	return
}
func (s *pbBuf) ReadInt32() (num uint32, err error) {
	err = readBytes(s.in, s.data_9[:4], 4)
	if err != nil {
		return
	}
	num = binary.BigEndian.Uint32(s.data_9)
	return
}

func (s *pbBuf) ReadInt16() (num uint16, err error) {
	err = readBytes(s.in, s.data_9[:2], 2)
	if err != nil {
		return
	}
	num = binary.BigEndian.Uint16(s.data_9)
	return
}
func (s *pbBuf) ReadInt8() (num uint8, err error) {
	err = readBytes(s.in, s.data_9[:1], 1)
	if err != nil {
		return
	}
	return uint8(s.data_9[0]), nil
}

func (s *pbBuf) WriteInt64(num uint64) error {
	binary.BigEndian.PutUint64(s.data_9, num)
	return writeBytes(s.out, s.data_9[:8])
}
func (s *pbBuf) WriteInt32(num uint32) error {
	binary.BigEndian.PutUint32(s.data_9, num)
	return writeBytes(s.out, s.data_9[:4])
}
func (s *pbBuf) WriteInt16(num uint16) error {
	binary.BigEndian.PutUint16(s.data_9, num)
	return writeBytes(s.out, s.data_9[:2])
}
func (s *pbBuf) WriteInt8(num uint8) error {
	return writeBytes(s.out, []byte{num})
}
func copyContentBuffer(dst io.Writer, src io.Reader) (err error) {
	r := NewReader(src)
	contentLen, _, err := r.ReadVarInt()
	if err != nil {
		return err
	}
	w := NewWriter(dst)
	w.WriteVarInt(contentLen)
	_, err = io.CopyN(dst, src, int64(contentLen))
	return
}
