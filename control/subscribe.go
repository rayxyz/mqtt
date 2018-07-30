package control

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mqtt/utils"
	"strings"
)

const (
	subScribeFixHeader    = 2
	subscribeVarHeaderLen = 2
)

// SubscribeHeader : Subscribe header
type SubscribeHeader struct {
	PackType  uint8
	RemainLen int
	PackID    uint16
}

// TopicFilter :
type TopicFilter struct {
	LenMSB     uint8
	LenLSB     uint8
	Filter     string
	RequestQoS uint8
}

// SubscribePayload :
type SubscribePayload struct {
	ClientID     string
	TopicFilters []*TopicFilter
}

// SubscribePacket : subscribe message packet
type SubscribePacket struct {
	Header  *SubscribeHeader
	Payload *SubscribePayload
}

// Marshal :
func (header *SubscribeHeader) marshal(payloadLen int) ([]byte, error) {
	if header == nil {
		return nil, errors.New("subscribe header is nil")
	}
	remainLen := subscribeVarHeaderLen + payloadLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("subscribe data is too large")
	}
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+subscribeVarHeaderLen)
	b[0] = 1<<7 | 1<<1
	b[1] = byte(remainLen)
	binary.BigEndian.PutUint16(b[2:4], header.PackID)
	if err := header.parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())

	return b, nil
}

// Parse : parse the subscribe packet
func (header *SubscribeHeader) parse(b []byte) error {
	header.PackType = b[0] >> 4
	header.PackID = binary.BigEndian.Uint16(b[2:4])
	remainLenDigits, err := ParseRemainLenDigits(b[1:5])
	if err != nil {
		return err
	}
	remainLen, err := DecodeRemainLen(remainLenDigits)
	if err != nil {
		return err
	}
	log.Println("remain_len_parsed => ", remainLen)
	header.RemainLen = remainLen

	return nil
}

// Marshal : marsal subscribe packet
func (p *SubscribePacket) Marshal() ([]byte, error) {
	var pbs []byte
	payload := p.Payload
	if strings.EqualFold(payload.ClientID, utils.Blank) {
		return nil, errors.New("invalid client identifier")
	}
	// for _, v := range payload.TopicFilters {
	//
	// }
	payloadBytes, err := json.Marshal(p.Payload)
	if err != nil {
		return nil, err
	}
	headerBytes, err := p.Header.marshal(len(payloadBytes))
	if err != nil {
		return nil, err
	}
	pbs = append(pbs, headerBytes...)
	pbs = append(pbs, payloadBytes...)

	return pbs, nil
}

// Parse : Parse subscribe packet
func (p *SubscribePacket) Parse(b []byte) error {
	headerLen := subScribeFixHeader + subscribeVarHeaderLen

	log.Println("subscribe header length => ", headerLen)
	header := new(SubscribeHeader)
	if err := header.parse(b); err != nil {
		return err
	}
	p.Header = header

	var payload SubscribePayload
	if err := json.Unmarshal(b[headerLen:], &payload); err != nil {
		return err
	}
	p.Payload = &payload

	return nil
}

func (header *SubscribeHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PackType=%d RemainLen=%d PackID=%d", header.PackType, header.RemainLen, header.PackID)
}
