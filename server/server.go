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
	if b[0]>>4 == control.CONNECT {
		fmt.Println("Connection packet detected.")
		packet, err := control.ParseConnectPacket(b)
		if !strings.EqualFold(store.ClientIDMap[packet.Payload.ClientID], utils.Blank) {
			log.Println("Repeated connection request, disconnecting....")
			conn.Close()
			delete(store.ClientIDMap, packet.Payload.ClientID)
			log.Println("Disconnected client => ", packet.Payload.ClientID)
			return
		}
		store.ClientIDMap[packet.Payload.ClientID] = packet.Payload.ClientID
		if err != nil {
			log.Println(err)
			return
		}
		conn.Write([]byte("header => " + (packet.Header.String()) + "  |||||| payload => " + fmt.Sprintf("%#v", packet.Payload) + "\n"))
	}
}
