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

// ConnectHeader : Connect header
type ConnectHeader struct {
	PackType         int
	RemainLen        int
	ProtocNameLenLSB int
	ProtocNameLenMSB int
	ProtocName       string
	ProtocLevel      int
	Flags            int
	KeepAliveLSB     int
	KeepAliveMSB     int
	KeepAlive        int
}

// ConnectPayload : Connection payload
type ConnectPayload struct {
	ClientID  string
	WillTopic string
	WillMsg   string
	UserName  string
	Password  string
}

// ConnectPacket : Connect packet
type ConnectPacket struct {
	Header  *ConnectHeader
	Payload *ConnectPayload
}

// Marshal : Marshal header to bytes
func (header *ConnectHeader) Marshal(payloadLen int, payload *ConnectPayload) ([]byte, error) {
	if header == nil {
		return nil, errors.New("connect header is nil")
	}
	remainLen := VarHeaderLen + payloadLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("connect data is too big")
	}
	remainLenBytes := GenRemainLenBytes(remainLen)
	fixedHeaderLen := len(remainLenBytes) + 1
	b := make([]byte, fixedHeaderLen+VarHeaderLen)
	b[0] = 1 << 4
	for i, v := range remainLenBytes {
		b[i+1] = v
	}
	b[fixedHeaderLen] = 0        // Length MSB
	b[fixedHeaderLen+1] = 1 << 2 // Length LSB
	// Protocol name
	b[fixedHeaderLen+2] = 'M'
	b[fixedHeaderLen+3] = 'Q'
	b[fixedHeaderLen+4] = 'T'
	b[fixedHeaderLen+5] = 'T'
	b[fixedHeaderLen+6] = 1 << 2 // Protoc level
	// Connect flags
	if !strings.EqualFold(payload.UserName, utils.Blank) {
		b[fixedHeaderLen+7] = 1 << 7 // Set user name flag to 1
	}
	if !strings.EqualFold(payload.Password, utils.Blank) {
		b[fixedHeaderLen+7] |= 1 << 6 // Set password flag to 1
	}
	b[fixedHeaderLen+7] &^= 1 << 5 // Set will retain flag to 0
	// Set QoS flags to 01
	b[fixedHeaderLen+7] &^= 1 << 4
	b[fixedHeaderLen+7] |= 1 << 3
	if !strings.EqualFold(payload.WillTopic, utils.Blank) {
		b[fixedHeaderLen+7] |= 1 << 2 // Set will flag to 1
	}
	if strings.EqualFold(payload.ClientID, utils.Blank) {
		b[fixedHeaderLen+7] |= 1 << 1 // Set Clean session to 1
	}
	b[fixedHeaderLen+7] &^= 1 << 0      // Clear reserved field
	b[fixedHeaderLen+8] = 0             // Keep Alive MSB
	b[fixedHeaderLen+9] = (1<<3 | 1<<1) // Keep Alive LSB
	if err := header.Parse(b); err != nil {
		log.Println(err)
	}
	fmt.Println("header.String() : \n", header.String())
	return b, nil
}

// Parse : parse connect header
func (header *ConnectHeader) Parse(b []byte) error {
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
	headerLen := 1 + len(remainLenDigits) + VarHeaderLen
	varHeaderStartIdx := 1 + len(remainLenDigits)
	header.RemainLen = remainLen
	header.ProtocNameLenLSB = int(b[varHeaderStartIdx])
	header.ProtocNameLenMSB = int(b[varHeaderStartIdx+1])
	header.ProtocName = string(b[varHeaderStartIdx+2 : varHeaderStartIdx+6])
	header.ProtocLevel = int(b[varHeaderStartIdx+6])
	header.Flags = int(b[varHeaderStartIdx+7])
	header.KeepAliveMSB = int(b[varHeaderStartIdx+8])
	log.Println("len(b) => ", len(b), "varHeaderStartIdx+9 => ", varHeaderStartIdx+9)
	header.KeepAliveLSB = int(b[varHeaderStartIdx+9])
	header.KeepAlive = int(binary.BigEndian.Uint16(b[varHeaderStartIdx+8 : headerLen]))
	fmt.Printf("parsed_header => %#v\n", header)
	return nil
}

// ParseConnectHeader : Parse connect header
func ParseConnectHeader(b []byte) (*ConnectHeader, error) {
	h := new(ConnectHeader)
	if err := h.Parse(b); err != nil {
		log.Println(err)
	}
	return h, nil
}

// ParseConnectPacket : Parse connect packet
func ParseConnectPacket(b []byte) (*ConnectPacket, error) {
	fixedHeaderLen, err := GetFixedHeaderLen(b[1:5])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	headerLen := fixedHeaderLen + VarHeaderLen
	log.Println("connect header length => ", headerLen)
	header, err := ParseConnectHeader(b[0:headerLen])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var payload ConnectPayload
	log.Println("length of bytes => ", len(b))
	if err = json.Unmarshal(b[headerLen-1:], &payload); err != nil {
		return nil, err
	}
	packet := &ConnectPacket{
		Header:  header,
		Payload: &payload,
	}
	return packet, nil
}

func (header *ConnectHeader) String() string {
	if header == nil {
		return "<nil>"
	}
	return fmt.Sprintf("remainlen=%d protoname=%s protolvl=%d flags=%d keepalive=%d", header.RemainLen, header.ProtocName, header.ProtocLevel, header.Flags, header.KeepAlive)
}
