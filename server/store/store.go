package store

import (
	"net"
)

var (
	// SessionMap : server sessions map to store server sessions
	SessionMap map[string]*Session
)

// Session : session storged in server
type Session struct {
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
	Conn                net.Conn
}
