package store

import (
	"net"
)

var (
	// ServerSessionMap : server sessions map to store server sessions
	ServerSessionMap map[string]*ServerSession
)

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
