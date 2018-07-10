package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mqtt/control"
	"mqtt/utils"
	"net"
	"time"
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
	payload := &control.ConnectPayload{
		ClientID:  utils.GenUUID(),
		WillTopic: "willtopic",
		WillMsg:   "will message fdsfsdfsdfsdfsdfsdfsdfsdfsdfsee",
		UserName:  "ray",
		Password:  "ray123",
	}
	connpack := &control.ConnectPacket{
		Header:  &control.ConnectHeader{},
		Payload: payload,
	}
	packdata, err := connpack.Marshal()
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("error of listening")
	}
	c.Conn = conn

	c.Conn.Write(packdata)

	connack := new(control.ConnAckPacket)
	datach := make(chan []byte)
	go func(datach chan []byte) {
		c.receive(datach)
	}(datach)
	connack.Parse(<-datach)
	fmt.Printf("ack => %#v\n", connack)
	if connack.Header.PackType != control.CONNACK {
		log.Println("wrong connect packet type")
		conn.Close()
	}

	log.Println("<<<<<<<<<<<<<<<<<< I will publish a message >>>>>>>>>>>>>>>>>>>>")

	c.Publish("Hello, my friend!!!")
}

func (c *Client) receive(ch chan []byte) {
	timeout := time.After(10 * time.Second)
loop:
	for {
		select {
		case <-timeout:
			log.Println("error: timeout")
			c.Conn.Close()
			break loop
		default:
			data, err := bufio.NewReader(c.Conn).ReadBytes('\n')
			if err != nil && err != io.EOF {
				log.Println("error of reading content")
				c.Conn.Close()
			}
			if len(data) > 0 {
				// log.Println("Before writing something to data channel.")
				// log.Println("data => ", data)
				ch <- data
				// log.Println("After writing something to data channel.")
				break loop
			}
		}
	}
}

// Publish message
func (c *Client) Publish(content interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(content); err != nil {
		log.Println(err)
		return
	}
	pubpack := &control.PublishPacket{
		Header: &control.PublishHeader{
			PacKID: 12345,
		},
		Payload: buf.Bytes(),
	}

	pbs, err := pubpack.Marshal()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.Conn.Write(pbs)
	if err != nil {
		log.Println(err)
	}
	log.Println("After sending the publish message.")

	puback := new(control.PublishAckPacket)
	datach := make(chan []byte)
	go func(datach chan []byte) {
		c.receive(datach)
	}(datach)
	ackPack, err := puback.Parse(<-datach)
	if err != nil {
		log.Println("parse publish ack pack error")
		return
	}
	log.Println("header.PacKID => ", pubpack.Header.PacKID, "ack.header.packid => ", ackPack.Header.PackID)
	if ackPack.Header.PackID != pubpack.Header.PacKID {
		log.Println("pack id does not matched")
		return
	}
	log.Println("publish ack packet => ", ackPack)
}
