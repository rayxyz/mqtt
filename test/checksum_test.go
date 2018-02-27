package test

import (
	"fmt"
	"testing"
)

func checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint16(0)
	for i := 0; i < csumcv; i += 2 {
		// s += uint32(b[i+1])<<8 | uint32(b[i])
		s += uint16(b[i])
		s += s & uint16(b[i+1])
	}
	// if csumcv&1 == 0 {
	// 	s += uint32(b[csumcv])
	// }
	// s = s>>16 + s&0xffff
	// s = s + s>>16

	return ^s
}

func TestChecksum(t *testing.T) {
	i := checksum([]byte{123, 2, 5, 6, 9, 10})
	fmt.Println("checksum => ", i)
	fmt.Println("reverse checksum => ", ^i)
	fmt.Printf("checksum bits => %08b\n", i)
	fmt.Printf("reverse checksum bits => %08b\n", ^i)
	flags := 0x00
	flags |= 0xff
	flags &= 0xfe
	fmt.Printf("flags %08b\n", flags)
	iflag := 122
	fmt.Printf("iflag => %08b\n", iflag)
	fmt.Printf("reverse iflag => %08b\n", ^iflag)
}
