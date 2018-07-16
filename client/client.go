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

// handle the packet from server
func (c *Client) handlePacket(data []byte) {
	cpt := data[0] >> 4

	log.Println("|||||||||||||||||||||||||||||| >>>>>>>>> cpt => ", cpt)

	switch cpt {
	case control.CONNACK:
		log.Println("ConnAck")
		c.handleConnAck(data)
	case control.PUBLISH:
		log.Println("<<Publish>>")
		c.handlePublish(data)
	case control.PUBACK:
		log.Println("PubAck")
		c.handlePublishAck(data)
	default:
		log.Println("no MQTT controll packet type matched")
		return
	}
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
}

func (c *Client) handleConnAck(data []byte) {
	log.Println("handling conn ack...")
}

func (c *Client) handlePublish(data []byte) {
	pubpack := new(control.PublishPacket)
	if err := pubpack.Parse(data); err != nil {
		log.Println("parse publish packet err => ", err)
	}
	log.Printf("publish pack from server => %v", pubpack)
	log.Println("payload of the publish packet from server => ", string(pubpack.Payload))
	log.Println("writing publish ack to server...")
	// ack := new(control.PublishAckPacket)
	// ack.Header = &control.PublishAckHeader{
	// 	PackID: pubpack.Header.PacKID,
	// }
	// ackpackData, err := ack.Marshal()
	// if err != nil {
	// 	log.Println("error of acking publish => ", err)
	// 	c.Conn.Close()
	// }
	// ackpackData = append(ackpackData, '\n')
	// c.Conn.Write(ackpackData)
}

func (c *Client) handlePublishAck(data []byte) {
	log.Println("handling publish ack...")
	puback := new(control.PublishAckPacket)
	ackPack, err := puback.Parse(data)
	if err != nil {
		log.Println("parse publish ack pack error")
		return
	}
	log.Println("header.PacKID => ", ackPack.Header.PackID, "ack.header.packid => ", ackPack.Header.PackID)
	// if ackPack.Header.PackID != pubpack.Header.PacKID {
	// 	log.Println("pack id does not matched")
	// 	return
	// }
	log.Println("publish ack packet => ", ackPack)
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
}
