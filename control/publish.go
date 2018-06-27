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
	PacKID    int
}

// PublishPayload : Publish payload
type PublishPayload struct {
	Content []byte
}

// PublishPacket : publish message packet
type PublishPacket struct {
	Header  *PublishHeader
	Payload *PublishPayload
}

// Marshal :
func (header *PublishHeader) Marshal(payloadLen int) ([]byte, error) {
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
	b[0] = 1 << 5
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

	if err := header.Parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())

	return b, nil
}

// Parse :
func (header *PublishHeader) Parse(b []byte) error {
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

// ParsePublishHeader : Parse pubish header
func ParsePublishHeader(b []byte) (*PublishHeader, error) {
	h := new(PublishHeader)
	if err := h.Parse(b); err != nil {
		log.Println(err)
	}
	return h, nil
}

// ParsePublishPacket : Parse publish packet
func ParsePublishPacket(b []byte) (*PublishPacket, error) {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	headerLen := fixedHeaderLen + publishVarHeaderLen
	log.Println("publish header length => ", headerLen)
	header, err := ParsePublishHeader(b[0:headerLen])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var payload PublishPayload
	log.Println("length of bytes => ", len(b))
	if err = json.Unmarshal(b[headerLen-1:], &payload); err != nil {
		return nil, err
	}
	packet := &PublishPacket{
		Header:  header,
		Payload: &payload,
	}
	return packet, nil
}

func (header *PublishHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d ", header.RemainLen)
}
