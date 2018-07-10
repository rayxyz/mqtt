package control

import (
	"errors"
	"fmt"
	"log"
)

const (
	varHeaderLen = 2
)

// ConnAckHeader : Connect acknowledgement header
type ConnAckHeader struct {
	PackType   int
	RemainLen  int
	Flags      int
	ReturnCode int
}

// ConnAckPacket : ConnAckPacket
type ConnAckPacket struct {
	Header *ConnAckHeader
}

// marshal header to bytes
func (header *ConnAckHeader) marshal() ([]byte, error) {
	if header == nil {
		return nil, errors.New("connect header is nil")
	}
	remainLen := varHeaderLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("connect data is too big")
	}
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+varHeaderLen)
	b[0] = 1 << 5
	for i, v := range remainLenBytes {
		b[i+1] = v
	}
	b[fixedHeaderLen] = 0
	b[fixedHeaderLen+1] = byte(header.ReturnCode)
	return b, nil
}

// parse connect ack header
func (header *ConnAckHeader) parse(b []byte) error {
	header.PackType = int(b[0] >> 4)
	remainLenDigits, err := ParseRemainLenDigits(b[1:3])
	if err != nil {
		return err
	}
	remainLen, err := DecodeRemainLen(remainLenDigits)
	if err != nil {
		return err
	}
	log.Println("remain_len_parsed => ", remainLen)
	header.RemainLen = remainLen
	varHeaderStartIdx := 1 + len(remainLenDigits)
	header.Flags = int(b[varHeaderStartIdx])
	header.ReturnCode = int(b[varHeaderStartIdx+1])
	fmt.Printf("parsed_header => %#v\n", header)
	return nil
}

// Marshal : marshal the connect ack packet
func (p *ConnAckPacket) Marshal() ([]byte, error) {
	var cpbs []byte
	headerBytes, err := p.Header.marshal()
	if err != nil {
		return nil, err
	}
	cpbs = append(cpbs, headerBytes...)
	return cpbs, nil
}

// Parse : Parse connect packet
func (p *ConnAckPacket) Parse(b []byte) error {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		return err
	}
	headerLen := fixedHeaderLen + varHeaderLen
	log.Println("connect acknowledge header length => ", headerLen)

	header := new(ConnAckHeader)
	if err := header.parse(b[0:headerLen]); err != nil {
		log.Println(err)
	}
	p.Header = header

	return nil
}

func (header *ConnAckHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d ", header.RemainLen)
}
