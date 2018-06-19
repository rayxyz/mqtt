package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mqtt/control"
	"net"
)

// Client : MQTT Client
type Client struct {
	Conn net.Conn
}

func main() {
	fmt.Println("connecting to the MQTT server")
	client := &Client{}
	client.Connect()
}

// Connect to the server
func (c *Client) Connect() {
	header := new(control.ConnectHeader)
	payload := &control.ConnectPayload{
		ClientID:  "KKKKK-XXX-YYY-ZZZ",
		WillTopic: "willtopic",
		WillMsg:   "will message fdsfsdfsdfsdfsdfsdfsdfsdfsdfsee",
		UserName:  "ray",
		Password:  "ray123",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic("marsh payload error")
	}
	headerBytes, err := header.Marshal(len(payloadBytes))
	if err != nil {
		panic("marshal header error")
	}
	var pack []byte
	pack = append(pack, headerBytes...)
	pack = append(pack, payloadBytes...)

	log.Println(pack)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("error of listening")
	}
	c.Conn = conn

	fmt.Fprintf(conn, string(pack))
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("error of reading content")
	}
	fmt.Println(data)
}
