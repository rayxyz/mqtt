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
	header := new(control.ConnectHeader)
	payload := &control.ConnectPayload{
		ClientID:  utils.GenUUID(),
		WillTopic: "willtopic",
		WillMsg:   "will message fdsfsdfsdfsdfsdfsdfsdfsdfsdfsee",
		UserName:  "ray",
		Password:  "ray123",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic("marsh payload error")
	}
	headerBytes, err := header.Marshal(len(payloadBytes), payload)
	if err != nil {
		panic("marshal header error")
	}
	var pack []byte
	pack = append(pack, headerBytes...)
	pack = append(pack, payloadBytes...)
	// pack = append(pack, '\n')

	log.Println(pack)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("error of listening")
	}
	c.Conn = conn

	// fmt.Fprintf(conn, string(pack))
	conn.Write(pack)

	log.Println("conn is nil after connected => ", conn == nil)

	// connack := new(control.ConnAckHeader)
	datach := make(chan int)
	c.receive(datach)
	data := <-datach
	log.Println("parsing received data... data => ", data)
	// connack.Parse(data)
	// fmt.Printf("ack => %#v\n", connack)
	// if connack.PackType != control.CONNACK {
	// 	log.Println("wrong connect packet type")
	// 	conn.Close()
	// }
	// loop:
	// 	for {
	// 		select {
	// 		case data := <-datach:
	// 			connack.Parse(data)
	// 			fmt.Printf("ack => %#v\n", connack)
	// 			if connack.PackType != control.CONNACK {
	// 				log.Println("wrong connect packet type")
	// 				conn.Close()
	// 			}
	// 			break loop
	// 		default:
	// 			break loop
	// 		}
	// 	}

	log.Println("<<<<<<<<<<<<<<<<<< I will publish a message >>>>>>>>>>>>>>>>>>>>")

	// c.Publish("Hello, my friend!!!")
}

func (c *Client) receive(ch chan int) {
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
				log.Println("Before writing something to data channel.")
				log.Println("data => ", data)
				ch <- 1
				log.Println("After writing something to data channel.")
			}
		}
	}
}

// Publish message
func (c *Client) Publish(content interface{}) {
	header := new(control.PublishHeader)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(content); err != nil {
		log.Println(err)
		return
	}
	payload := &control.PublishPayload{
		Content: buf.Bytes(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		panic("marsh payload error")
	}
	headerBytes, err := header.Marshal(len(payloadBytes))
	if err != nil {
		log.Println(err)
		panic("marshal header error")
	}

	var pack []byte
	pack = append(pack, headerBytes...)
	pack = append(pack, payloadBytes...)
	// pack = append(pack, '\n')

	log.Println(pack)

	// n, err := fmt.Fprintf(c.Conn, string(pack))
	n, err := c.Conn.Write(pack)
	if err != nil {
		log.Println(err)
	}
	log.Println("n => ", n)

	log.Println("After sending the publish message.")

	puback := new(control.PublishAckPacket)
	datach := make(chan []byte)
	// c.receive(datach)
loop:
	for {
		select {
		case data := <-datach:
			puback.Parse(data)
			log.Println("puback => ", puback)
			break loop
		default:
			break loop
		}
	}
}
