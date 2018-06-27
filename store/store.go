package store

import (
	"net"
)

var (
	// ClientIDMap : Client identifier map to store client IDs.
	ClientIDMap map[string]string
	// ServerSessionMap : server session map
	ServerSessionMap map[string]*ServerSession
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

// ServerSession : session storged in server
type ServerSession struct {
	ID                  string
	ClientID            string
	Type                int
	ConnectReceived     bool
	ConnAckSent         bool
	PublishReceived     bool
	PubAckSent          bool
	PubRecSent          bool
	PubRelReceived      bool
	PubComSent          bool
	SubscribeReceived   bool
	SubAckSent          bool
	UnsubscribeReceived bool
	UnsubAckSent        bool
	PingReqReceived     bool
	PingRespSent        bool
	DisconnectReceived  bool
	Connection          net.Conn
}

// ClientMetadata : storage the metadata of client
type ClientMetadata struct {
}

// ServerMetadata : storage the metadata of the server
type ServerMetadata struct {
}
