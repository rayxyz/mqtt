package client

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
	ID     string
	Datach chan []byte
	Conn   net.Conn
}

// Run the client
func (c *Client) Run(serverURL string) {
	fmt.Println("Running MQTT client...")
	// client := &Client{}
	ch := make(chan int)
	go c.Connect(ch, serverURL)
	<-ch
	// The goroutines block on sending to the unbuffered channel.
	// A minimal change unblocks the goroutines is to create
	// a buffered channel with capacity
	c.Datach = make(chan []byte, 1)

	ech := make(chan error)
	go func(datach chan []byte, ech chan error) {
		counter := 0
		// Since the buffer is not persistent across iterations, any messages
		// received before the new Reader is created will be lost. So, here I
		// create the reader outside of the for loop.
		reader := bufio.NewReader(c.Conn)
		for {
			data, err := reader.ReadBytes('\r')
			if err != nil && err != io.EOF {
				log.Println("error of reading content")
				c.Conn.Close()
				ech <- err
				break
			}
			if len(data) > 0 {
				counter++
				log.Println("counter => ", counter, " len(data) => ", len(data))
				datach <- data
			}
		}
	}(c.Datach, ech)

	for {
		select {
		case data := <-c.Datach:
			log.Println("data >>>>>> received => ", data)
			go c.handlePacket(data)
		case err := <-ech:
			log.Println(err)
			break
		default:
			// do nothing here
		}
	}
}

// handle the packet from server
func (c *Client) handlePacket(data []byte) {
	cpt := data[0] >> 4

	log.Println("|||||||||||||||||||||||||||||| >>>>>>>>> cpt => ", cpt)

	switch cpt {
	case control.CONNACK:
		log.Println("<<ConnAck>>")
		c.handleConnAck(data)
	case control.PUBLISH:
		log.Println("<<Publish>>")
		c.handlePublish(data)
	case control.PUBACK:
		log.Println("<<PubAck>>")
		c.handlePublishAck(data)
	default:
		log.Println("no MQTT controll packet type matched")
		return
	}
}

// Connect to the server
func (c *Client) Connect(ch chan int, addr string) {
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
	packdata = append(packdata, '\r')

	conn, err := net.Dial("tcp", addr)
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
	pbs = append(pbs, '\r')

	_, err = c.Conn.Write(pbs)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("After sending the publish message.")
}

// Subscribe topics
func (c *Client) Subscribe() {
	log.Println("subscribing...")
	subpack := &control.SubscribePacket{
		Header: &control.SubscribeHeader{
			PacKID: 12345,
		},
		Payload: &control.SubscribePayload{
			ClientID: c.ID,
			TopicFilters: []*control.TopicFilter{
				&control.TopicFilter{
					Filter:     "/status/*",
					RequestQoS: 0,
				},
			},
		},
	}

	subpackBytes, err := subpack.Marshal()
	if err != nil {
		log.Println(err)
		return
	}
	subpackBytes = append(subpackBytes, '\r')

	_, err = c.Conn.Write(subpackBytes)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("subscribe data has been sent")
}
