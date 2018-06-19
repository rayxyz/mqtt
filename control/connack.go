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
	Header *ConnectHeader
}

// Marshal : Marshal header to bytes
func (header *ConnAckHeader) Marshal(returnCode int) ([]byte, error) {
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
	b[fixedHeaderLen+1] = byte(returnCode)
	if err := header.Parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())
	return b, nil
}

// Parse : parse connect header
func (header *ConnAckHeader) Parse(b []byte) error {
	header.PackType = int(b[0] >> 5)
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

// ParseConnAckHeader : Parse connect header
func ParseConnAckHeader(b []byte) (*ConnectHeader, error) {
	h := new(ConnectHeader)
	if err := h.Parse(b); err != nil {
		log.Println(err)
	}
	return h, nil
}

// ParseConnAckPacket : Parse connect packet
func ParseConnAckPacket(b []byte) (*ConnAckPacket, error) {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	headerLen := fixedHeaderLen + varHeaderLen
	log.Println("connect acknowledge header length => ", headerLen)
	header, err := ParseConnectHeader(b[0:headerLen])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	packet := &ConnAckPacket{
		Header: header,
	}
	return packet, nil
}

func (header *ConnAckHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d ", header.RemainLen)
}
