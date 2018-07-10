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

type mqttServer struct{}

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

func handleConn(conn net.Conn) {
	ch := make(chan []byte)
	ech := make(chan error)
	//
	// buf := make([]byte, 0, 4096)
	// packLenAlereadyParsed := false
	// packLen := 0

	go func(ch chan []byte, ech chan error) {
		tmp := make([]byte, 1024)
		counter := 0
		for {
			n, err := conn.Read(tmp)
			if err != nil && err != io.EOF {
				ech <- err
				break
			}
			counter++
			if n > 0 {
				log.Println("counter => ", counter, " n => ", n)
				if n < len(tmp) {
					tmp = tmp[:n]
				}
				ch <- tmp
			}
		}
	}(ch, ech)

	for {
		select {
		case data := <-ch:
			log.Println("data => ", data)
			// if !packLenAlereadyParsed {
			// 	packLenParsed, err := control.GetPackLen(data[1:5])
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	log.Println("pack_len => ", packLenParsed)
			// 	// packLen = packLenParsed
			// 	packLenAlereadyParsed = true
			// }
			// buf = append(buf, data[:len(data)]...)
			// if len(buf) >= packLen {
			// 	server.handlePacket(conn, buf)
			// 	buf = append(buf[:0])
			// }
			//
			// packLenParsed, err := control.GetPackLen(data[1:5])
			// if err != nil {
			// 	log.Println(err)
			// }
			// log.Println("pack_len => ", packLenParsed)
			server.handlePacket(conn, data)
		case err := <-ech:
			log.Println(err)
			break
		}
	}
}

// handlePacket : Read MQTT packet from stream and handle it
func (s *mqttServer) handlePacket(conn net.Conn, b []byte) {
	cpt := b[0] >> 4

	log.Println("|||||||||||||||||||||||||||||| >>>>>>>>> cpt => ", cpt)

	switch cpt {
	case control.CONNECT:
		s.handleConnect(conn, b)
	case control.PUBLISH:
		fmt.Println("<<<Publish>>>")
		s.handlePublish(conn, b)
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

func (s *mqttServer) handleConnect(conn net.Conn, data []byte) {
	connpack := new(control.ConnectPacket)
	if err := connpack.Parse(data); err != nil {
		log.Println(err)
		return
	}
	log.Println("protocol name => ", connpack.Header.ProtocName)
	ackpack := new(control.ConnAckPacket)
	if !strings.EqualFold(connpack.Header.ProtocName, utils.ProtocName) {
		ackpack.Header = &control.ConnAckHeader{
			ReturnCode: 1,
		}
		ackpackData, err := ackpack.Marshal()
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		packdata := append(ackpackData, '\n')
		conn.Write(packdata)
		conn.Close()
		return
	}
	if !strings.EqualFold(connpack.Payload.ClientID, utils.Blank) {
		if !strings.EqualFold(store.ClientIDMap[connpack.Payload.ClientID], utils.Blank) {
			ackpack.Header = &control.ConnAckHeader{
				ReturnCode: 2,
			}
			ackpackData, err := ackpack.Marshal()
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			packdata := append(ackpackData, '\n')
			conn.Write(packdata)
			conn.Close()
			delete(store.ClientIDMap, connpack.Payload.ClientID)
			log.Println("disconnected client => ", connpack.Payload.ClientID)
			return
		}
	} else {
		// clean session not set
		if connpack.Header.Flags&(1<<1) == 0 {
			ackpack.Header = &control.ConnAckHeader{
				ReturnCode: 2,
			}
			ackpackData, err := ackpack.Marshal()
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			packdata := append(ackpackData, '\n')
			conn.Write(packdata)
			conn.Close()
			return
		}
	}
	store.ClientIDMap[connpack.Payload.ClientID] = connpack.Payload.ClientID
	ackpack.Header = &control.ConnAckHeader{
		ReturnCode: 0,
	}
	ackpackData, err := ackpack.Marshal()
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	packdata := append(ackpackData, '\n')
	conn.Write(packdata)
}

func (s *mqttServer) handlePublish(conn net.Conn, data []byte) {
	pubpack := new(control.PublishPacket)
	if err := pubpack.Parse(data); err != nil {
		log.Println("parse publish packet err => ", err)
	}
	log.Printf("publish pack => %v", pubpack)
	log.Println("payload => ", string(pubpack.Payload))
	log.Println("writing publish ack...")
	ack := new(control.PublishAckPacket)
	ack.Header = &control.PublishAckHeader{
		PackID: 12345,
	}
	packBytes, err := ack.Marshal()
	if err != nil {
		log.Println("error of acking publish => ", err)
		conn.Close()
	}
	packBytes = append(packBytes, '\n')
	conn.Write(packBytes)
}
