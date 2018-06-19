package control

import (
	"errors"
	"io"
	"log"
	"mqtt/utils"
	"net"
	"time"
)

// Constants
const (
	VarHeaderLen = 10
	// 4 bytes to represent data in variable header and
	// the payload, the top bit excluded.
	MaxRemainLen = 1 << 28
)

// Header : MQTT header interface
type Header interface {
	Marshal([]byte)
	Parse() []byte
	String() string
}

// EncodeRemainLen : Encode remaining length value
func EncodeRemainLen(remainLen int) []int {
	var digits []int
	for {
		digit := remainLen % 128
		remainLen = remainLen / 128
		if remainLen > 0 {
			digit |= 128
		}
		digits = append(digits, digit)
		if remainLen <= 0 {
			break
		}
	}
	return digits
}

// DecodeRemainLen : Decode remaining length digits to value
func DecodeRemainLen(digits []int) (int, error) {
	value := 0
	multiplier := 1
	for _, v := range digits {
		value += (v & 127) * multiplier
		if multiplier > 128*128*128 {
			return 0, errors.New("malformed remain length")
		}
		multiplier *= 128
	}
	return value, nil
}

// ParseRemainLenDigits : Parse remaining digits from header
func ParseRemainLenDigits(evalBytes []byte) ([]int, error) {
	if evalBytes == nil || len(evalBytes) < 1 {
		return nil, errors.New("invalid evaluating bytes")
	}
	var remainLenDigits []int
	for i := 0; i < len(evalBytes); i++ {
		remainLenDigits = append(remainLenDigits, int(evalBytes[i]))
		// If the top bit of is 1, then continue to get remaining length field bit
		if evalBytes[i]>>7 == 1 {
			continue
		} else {
			break
		}
	}
	return remainLenDigits, nil
}

// GetFixedHeaderLen : Get fixed header length
func GetFixedHeaderLen(evalBytes []byte) (int, error) {
	digits, err := ParseRemainLenDigits(evalBytes)
	if err != nil {
		return 0, err
	}
	return len(digits) + 1, nil
}

// GetRemainLen : Get remaining length
func GetRemainLen(evalBytes []byte) (int, error) {
	remainLenDigits, err := ParseRemainLenDigits(evalBytes)
	if err != nil {
		return 0, err
	}
	remainLen, err := DecodeRemainLen(remainLenDigits)
	if err != nil {
		return 0, err
	}
	return remainLen, nil
}

// GenRemainLenBytes : generate remaining length bytes
func GenRemainLenBytes(remainLen int) []byte {
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
	return remainLenBytes
}

// GetPackLen : Get the packet length
func GetPackLen(evalBytes []byte) (int, error) {
	remainLenDigits, err := ParseRemainLenDigits(evalBytes)
	if err != nil {
		return 0, errors.New("parse remaining length digits err")
	}
	remainLen, err := DecodeRemainLen(remainLenDigits)
	if err != nil {
		return 0, errors.New("get protocal header reamain length error")
	}
	log.Println("remaining length => ", remainLen)
	packLen := 1 + len(remainLenDigits) + remainLen
	log.Println("packet length => ", packLen)
	return packLen, nil
}

func calculateRemainLenBytes(remainLen int) (int, error) {
	remainBytesLen := 1
	// calculate how many bytes should be assigned to the remaining lenth field
	if remainLen < 1<<7 {
		remainBytesLen = 1
	} else if remainLen < 1<<14 {
		remainBytesLen = 2
	} else if remainLen < 1<<21 {
		remainBytesLen = 3
	} else if remainLen < 1<<28 {
		remainBytesLen = 4
	} else {
		return 0, errors.New("calculate remain bytes error")
	}

	return remainBytesLen, nil
}

// ReadPacket : Read MQTT packet from stream
func ReadPacket(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 128)
	packLenAlereadyParsed := false
	packLen := 0
	timeout := time.After(utils.ReadPackTimeout * time.Second)
loop:
	for {
		select {
		case <-timeout:
			return nil, errors.New("read TCP packet timeout")
		default:
			n, err := conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					break loop
				}
				log.Println("read data error")
				return nil, err
			}
			if !packLenAlereadyParsed {
				packLenParsed, err := GetPackLen(tmp[1:5])
				if err != nil {
					log.Println(err)
					return nil, err
				}
				packLen = packLenParsed
				packLenAlereadyParsed = true
			}
			buf = append(buf, tmp[:n]...)
			if len(buf) >= packLen {
				break loop
			}
		}
	}
	return buf, nil
}
