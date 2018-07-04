package control

import (
	"errors"
	"fmt"
	"log"
)

const (
	publishAckVarHeaderLen = 2
)

// PublishAckHeader : publish acknowledgement header
type PublishAckHeader struct {
	PackType   int
	RemainLen  int
	Flags      int
	ReturnCode int
}

// PublishAckPacket : PublishAckPacket
type PublishAckPacket struct {
	Header *PublishAckHeader
}

// Marshal : Marshal header to bytes
func (header *PublishAckHeader) Marshal(returnCode int) ([]byte, error) {
	if header == nil {
		return nil, errors.New("publish header is nil")
	}
	remainLen := publishAckVarHeaderLen
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+publishAckVarHeaderLen)
	b[0] = 1 << 6
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
func (header *PublishAckHeader) Parse(b []byte) error {
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

// ParsePublishectHeader : Parse connect header
func ParsePublishectHeader(b []byte) (*PublishAckHeader, error) {
	h := new(PublishAckHeader)
	if err := h.Parse(b); err != nil {
		log.Println(err)
	}
	return h, nil
}

// ParsePublishAckPacket : Parse connect packet
func ParsePublishAckPacket(b []byte) (*PublishAckPacket, error) {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	headerLen := fixedHeaderLen + varHeaderLen
	log.Println("connect acknowledge header length => ", headerLen)
	header, err := ParsePublishectHeader(b[0:headerLen])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	packet := &PublishAckPacket{
		Header: header,
	}
	return packet, nil
}

func (header *PublishAckHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d ", header.RemainLen)
}
