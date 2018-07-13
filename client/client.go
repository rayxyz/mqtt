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
)

// Client : MQTT Client
type Client struct {
	Datach chan []byte
	Conn   net.Conn
}

func main() {
	fmt.Println("connecting to the MQTT server")
	client := &Client{}
	ch := make(chan int)
	go client.Connect(ch)
	<-ch
	// The goroutines block on sending to the unbuffered channel.
	// A minimal change unblocks the goroutines is to create
	// a buffered channel with capacity
	client.Datach = make(chan []byte, 1)
	go client.receive()
	i := 0
	for {
		select {
		case data := <-client.Datach:
			log.Println("data >>>>>> received => ", data)
			if i == 0 {
				client.Publish("Hello, my friend!!!")
				// time.Sleep(2 * time.Second)
				i++
			}
			client.handlePacket(data)
		default:
			// do nothing here
		}
	}

}

// handle the packet from server
func (c *Client) handlePacket(packet []byte) {
	cpt := packet[0] >> 4

	log.Println("|||||||||||||||||||||||||||||| >>>>>>>>> cpt => ", cpt)

	switch cpt {
	case control.CONNACK:
		c.handleConnAck(packet)
	case control.PUBACK:
		c.handlePublishAck(packet)
	default:
		log.Println("no MQTT controll packet type matched")
		return
	}
}

func (c *Client) handleConnAck(data []byte) {
	log.Println("handling conn ack...")
}

func (c *Client) handlePublishAck(packet []byte) {
	log.Println("handling publish ack...")
}

// Connect to the server
func (c *Client) Connect(ch chan int) {
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

	ch <- 1

	c.Conn.Write(packdata)

	// connack := new(control.ConnAckPacket)
	// // datach := make(chan []byte)
	// // go func(datach chan []byte) {
	// // 	c.receive(datach)
	// // }(datach)
	// connack.Parse(<-c.Datach)
	// fmt.Printf("ack => %#v\n", connack)
	// if connack.Header.PackType != control.CONNACK {
	// 	log.Println("wrong connect packet type")
	// 	conn.Close()
	// }
}

func (c *Client) receive() {
	for {
		data, err := bufio.NewReader(c.Conn).ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Println("error of reading content")
			c.Conn.Close()
		}
		if len(data) > 0 {
			log.Println("Before writing something to data channel.")
			c.Datach <- data
			log.Println("After writing something to data channel.")
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

	// puback := new(control.PublishAckPacket)
	// // datach := make(chan []byte)
	// // go func(datach chan []byte) {
	// // 	c.receive(datach)
	// // }(datach)
	// ackPack, err := puback.Parse(<-c.Datach)
	// if err != nil {
	// 	log.Println("parse publish ack pack error")
	// 	return
	// }
	// log.Println("header.PacKID => ", pubpack.Header.PacKID, "ack.header.packid => ", ackPack.Header.PackID)
	// if ackPack.Header.PackID != pubpack.Header.PacKID {
	// 	log.Println("pack id does not matched")
	// 	return
	// }
	// log.Println("publish ack packet => ", ackPack)
}
