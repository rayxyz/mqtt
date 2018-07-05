package control

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"mqtt/utils"
)

const (
	publishAckFixedHeaderLen = 2
	publishAckVarHeaderLen   = 2
)

// PublishAckHeader : publish acknowledgement header
type PublishAckHeader struct {
	PackType  int
	RemainLen int
	PackID    uint16
}

// PublishAckPacket : PublishAckPacket
type PublishAckPacket struct {
	Header *PublishAckHeader
}

// Marshal : Marshal header to bytes
func (header *PublishAckHeader) Marshal() ([]byte, error) {
	if header == nil {
		return nil, errors.New("publish header is nil")
	}
	b := make([]byte, publishAckFixedHeaderLen+publishAckVarHeaderLen)
	b[0] = 1<<6 | 1<<5 | 1<<4
	b[1] = 1 << 1
	log.Println("b => ", b)
	binary.BigEndian.PutUint16(b[2:4], header.PackID)
	log.Println(">>>>>>=============<<<<<<<<")
	if err := header.Parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())
	return b, nil
}

// Parse : parse connect header
func (header *PublishAckHeader) Parse(b []byte) error {
	header.PackType = int(b[0] >> 4)
	header.RemainLen = int(b[1] >> 1)
	header.PackID = binary.BigEndian.Uint16(b[2:4])
	fmt.Printf("parsed_header => %#v\n", header)
	return nil
}

// Marshal : marshal the publish ack packet
func (p *PublishAckPacket) Marshal() ([]byte, error) {
	return p.Header.Marshal()
}

// ParseHeader : Parse publish acknowledgement header
func (p *PublishAckPacket) ParseHeader(b []byte) (*PublishAckHeader, error) {
	h := new(PublishAckHeader)
	if err := h.Parse(b); err != nil {
		log.Println(err)
	}
	return h, nil
}

// Parse : Parse connect packet
func (p *PublishAckPacket) Parse(b []byte) (*PublishAckPacket, error) {
	header, err := p.ParseHeader(b[0 : publishAckFixedHeaderLen+publishAckVarHeaderLen])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	packet := &PublishAckPacket{
		Header: header,
	}
	log.Println("publish ack packet => ", packet)
	return packet, nil
}

func (header *PublishAckHeader) String() string {
	if header == nil {
		return utils.Nil
	}
	return fmt.Sprintf("PackType=%d RemainLen=%d PackID=%d",
		header.PackType, header.RemainLen, header.PackID)
}
