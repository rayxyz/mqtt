package control

import (
	"errors"
	"fmt"
	"log"
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
type PublishPayload []byte

// Marshal :
func (header *PublishHeader) Marshal(payloadLen int, payload *PublishPayload) ([]byte, error) {
	if header == nil {
		return nil, errors.New("publish header is nil")
	}
	remainLen := VarHeaderLen + payloadLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("publish data is too large")
	}
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+VarHeaderLen)
	b[0] = 1 << 4
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

func (header *PublishHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d ", header.RemainLen)
}
