package control

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
func (header *ConnectHeader) Marshal(payloadLen int) ([]byte, error) {
	if header == nil {
		return nil, errors.New("connect header is nil")
	}
	// payLoadLen := len(payLoad.ClientID) + len(payLoad.WillTopic) + len(payLoad.WillMsg) +
	// 	len(payLoad.UserName) + len(payLoad.Password)
	remainLen := VarHeaderLen + payloadLen
	if remainLen > MaxRemainLen {
		return nil, errors.New("connect data is too big")
	}
	digits := EncodeRemainLen(remainLen)
	remainLenBytes := make([]byte, len(digits))
	if len(digits) == 1 {
		remainLenBytes[0] = byte(remainLen)
	} else if len(digits) == 2 {
		// binary.BigEndian.PutUint16(remainLenBytes[0:len(digits)], uint16(remainLen))
		remainLenBytes[0] = byte(digits[0])
		remainLenBytes[1] = byte(digits[1])
	} else if len(digits) == 3 {
		// binary.BigEndian.PutUint32(remainLenBytes[0:len(digits)], uint32(remainLen))
		remainLenBytes[0] = byte(digits[0])
		remainLenBytes[1] = byte(digits[1])
		remainLenBytes[2] = byte(digits[2])
	} else {
		remainLenBytes[0] = byte(digits[0])
		remainLenBytes[1] = byte(digits[1])
		remainLenBytes[2] = byte(digits[2])
		remainLenBytes[3] = byte(digits[3])
	}
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
	b[fixedHeaderLen+7] = 1 << 7   // Set user name flag to 1
	b[fixedHeaderLen+7] |= 1 << 6  // Set password flag to 1
	b[fixedHeaderLen+7] &^= 1 << 5 // Set will retain flag to 0
	// Set QoS flags to 01
	b[fixedHeaderLen+7] &^= 1 << 4
	b[fixedHeaderLen+7] |= 1 << 3
	b[fixedHeaderLen+7] |= 1 << 2       // Set will flag to 1
	b[fixedHeaderLen+7] |= 1 << 1       // Set Clean session to 1
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
