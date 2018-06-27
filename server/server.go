package main

import (
	"fmt"
	"io"
	"log"
	"mqtt/control"
	"mqtt/store"
	"mqtt/utils"
	"net"
	"strings"
)

type mqttServer struct {
}

var server = new(mqttServer)

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

// ReadPacket : Read MQTT packet from stream
func handleConn(conn net.Conn) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 128)
	packLenAlereadyParsed := false
	packLen := 0

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				// buf = make([]byte, 0, 4096)
				// tmp = make([]byte, 128)
				// packLenAlereadyParsed = false
				// packLen = 0
			} else {
				log.Println("read data error")
			}
		}
		if !packLenAlereadyParsed {
			packLenParsed, err := control.GetPackLen(tmp[1:5])
			if err != nil {
				log.Println(err)
			}
			packLen = packLenParsed
			packLenAlereadyParsed = true
		}
		buf = append(buf, tmp[:n]...)
		if len(buf) >= packLen {
			go server.handlePacket(conn, buf)
		}
	}
}

func (s *mqttServer) handlePacket(conn net.Conn, b []byte) {
	cpt := b[0] >> 4

	log.Println("|||||||||||||||||||||||||||||| >>>>>>>>> cpt => ", cpt)

	switch cpt {
	case control.CONNECT:
		packet, err := control.ParseConnectPacket(b)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("<<<Connect>>> client_id => ", packet.Payload.ClientID)
		s.handleConnect(conn, packet)
	case control.PUBLISH:
		fmt.Println("<<<Publish>>>")
		packet, err := control.ParsePublishPacket(b)
		if err != nil {
			log.Println(err)
			return
		}
		s.handlePublish(conn, packet)
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

func (s *mqttServer) handleConnect(conn net.Conn, packet *control.ConnectPacket) {
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

func (s *mqttServer) handlePublish(conn net.Conn, packet *control.PublishPacket) {
	log.Println("publish pack => ", packet)
}
