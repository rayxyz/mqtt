package main

import (
	"fmt"
	"log"
	"mqtt/control"
	"mqtt/store"
	"mqtt/utils"
	"net"
	"strings"
)

func main() {
	fmt.Println("I am the primitive MQTT server, and I am alive...")
	fmt.Println("Server running on port => ", 8080)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic("error of creating mqtt server")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error of accept connection")
		}
		go handleConn(conn)
	}
}

func init() {
	fmt.Println("I am fucking init....")
	store.ClientIDMap = make(map[string]string, 10)
}

func handleConn(conn net.Conn) {
	// defer conn.Close()
	b, err := control.ReadPacket(conn)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(" data lenth => ", len(b), "data received => ", b)
	packet, err := control.ParseConnectPacket(b)
	if err != nil {
		log.Println(err)
		return
	}
	switch packet.Header.PackType {
	case control.CONNECT:
		fmt.Println("<<<Connect>>>")
		handleConnect(conn, packet)
	case control.PUBLISH:
		fmt.Println("<<<Publish>>>")
	case control.SUBSCRIBE:
		fmt.Println("<<<Subscribe>>>")
	case control.UNSUBSCRIBE:
		fmt.Println("<<<Unsubscribe>>>")
	case control.PINGREQ:
		fmt.Println("<<<PingReq>>>")
	case control.DISCONNECT:
		fmt.Println("<<<Disconnect>>>")
	default:
		log.Println("no MQTT controll packet type matched")
		conn.Close()
		return
	}
}

func handleConnect(conn net.Conn, packet *control.ConnectPacket) {
	log.Println("protocol name => ", packet.Header.ProtocName)
	ackHeader := new(control.ConnAckHeader)
	if !strings.EqualFold(packet.Header.ProtocName, utils.ProtocName) {
		headerBytes, _ := ackHeader.Marshal(1)
		headerBytes = append(headerBytes, '\n')
		conn.Write(headerBytes)
		conn.Close()
		return
	}
	if !strings.EqualFold(store.ClientIDMap[packet.Payload.ClientID], utils.Blank) {
		headerBytes, _ := ackHeader.Marshal(2)
		headerBytes = append(headerBytes, '\n')
		conn.Write(headerBytes)
		conn.Close()
		delete(store.ClientIDMap, packet.Payload.ClientID)
		log.Println("disconnected client => ", packet.Payload.ClientID)
		return
	} else {
		// clean session not set
		if packet.Header.Flags&(1<<1) == 0 {
			headerBytes, _ := ackHeader.Marshal(2)
			headerBytes = append(headerBytes, '\n')
			conn.Write(headerBytes)
			conn.Close()
			return
		}
	}
	store.ClientIDMap[packet.Payload.ClientID] = packet.Payload.ClientID
	ackHeaderBytes, _ := ackHeader.Marshal(0)
	ackHeaderBytes = append(ackHeaderBytes, '\n')
	conn.Write(ackHeaderBytes)
}
