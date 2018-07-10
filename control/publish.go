package control

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

const (
	publishVarHeaderLen = 7
)

// PublishHeader : Publish header
type PublishHeader struct {
	PackType  int
	DUPFlag   int
	QoSLevel  int
	RETAIN    int
	RemainLen int

	TopicName []byte
	PacKID    uint16
}

// PublishPacket : publish message packet
type PublishPacket struct {
	Header  *PublishHeader
	Payload []byte
}

// Marshal :
func (header *PublishHeader) marshal(payloadLen int) ([]byte, error) {
	if header == nil {
		return nil, errors.New("publish header is nil")
	}
	remainLen := publishVarHeaderLen + payloadLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("publish data is too large")
	}
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+publishVarHeaderLen)
	b[0] = 1<<5 | 1<<4
	if header.DUPFlag == 1 {
		b[0] |= 1 << 3
	}
	if header.QoSLevel == 0 {
		// do nothing here
	} else if header.QoSLevel == 1 {
		b[0] |= 1 << 1
	} else if header.QoSLevel == 2 {
		b[0] |= 1 << 2
	} else {
		return nil, errors.New("invalid QoS level")
	}
	for i, v := range remainLenBytes {
		b[i+1] = v
	}

	if err := header.parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())

	return b, nil
}

// Parse : parse the publish packet
func (header *PublishHeader) parse(b []byte) error {
	header.PackType = int(b[0] >> 4)
	remainLenDigits, err := ParseRemainLenDigits(b[1:5])
	if err != nil {
		return err
	}
	remainLen, err := DecodeRemainLen(remainLenDigits)
	if err != nil {
		return err
	}
	log.Println("remain_len_parsed => ", remainLen)

	return nil
}

// Marshal : marsal publish packet
func (p *PublishPacket) Marshal() ([]byte, error) {
	var pbs []byte
	headerBytes, err := p.Header.marshal(len(p.Payload))
	if err != nil {
		return nil, err
	}
	pbs = append(pbs, headerBytes...)

	payloadBytes, err := json.Marshal(p.Payload)
	if err != nil {
		return nil, err
	}
	pbs = append(pbs, payloadBytes...)

	return pbs, nil
}

// Parse : Parse publish packet
func (p *PublishPacket) Parse(b []byte) error {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		return err
	}
	headerLen := fixedHeaderLen + publishVarHeaderLen
	log.Println("publish header length => ", headerLen)

	header := new(PublishHeader)
	if err = header.parse(b); err != nil {
		return err
	}
	p.Header = header

	var payload []byte
	log.Println("length of bytes => ", len(b))
	if err = json.Unmarshal(b[headerLen:], &payload); err != nil {
		return err
	}
	p.Payload = payload

	return nil
}

func (header *PublishHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PackType=%d RemainLen=%d PackID=%d", header.PackType, header.RemainLen, header.PacKID)
}
