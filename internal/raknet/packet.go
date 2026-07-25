package raknet

import (
	"bytes"
	"encoding/binary"
)

type Packet struct {
	Buffer *bytes.Buffer
}

func NewPacket(data []byte) *Packet {
	return &Packet{
		Buffer: bytes.NewBuffer(data),
	}
}

func (p *Packet) ReadByte() (byte, error) {
	return p.Buffer.ReadByte()
}

func (p *Packet) ReadUint16() (uint16, error) {
	var val uint16
	err := binary.Read(p.Buffer, binary.BigEndian, &val)
	return val, err
}

func (p *Packet) ReadUint32() (uint32, error) {
	var val uint32
	err := binary.Read(p.Buffer, binary.BigEndian, &val)
	return val, err
}
