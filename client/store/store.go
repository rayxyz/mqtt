package store

import (
	"net"
)

var (
	// ClientSessionMap : client session map
	ClientSessionMap map[string]*ClientSession
)

// ClientSession : session stored in client
type ClientSession struct {
	ID               string
	ClientID         string
	Type             int
	ConnectSent      bool
	ConnAckReceived  bool
	PublishSent      bool
	PubAckReceived   bool
	PubRecRecived    bool
	PubRelSent       bool
	PubComReceived   bool
	SubscribeSent    bool
	SubAckReceived   bool
	UnsubscribeSent  bool
	UnsubAckReceived bool
	PingReqSent      bool
	PingRespReceived bool
	Connection       net.Conn
	PackID           uint16
}
