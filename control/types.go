package control

const (
	_ = iota
	// CONNECT : Connect
	CONNECT
	// CONNACK : Connect acknowledgement
	CONNACK
	// PUBLISH : Publish
	PUBLISH
	// PUBACK : Publish acknowledgement
	PUBACK
	// PUBREC : Publish receive
	PUBREC
	// PUBREL : Publish release
	PUBREL
	// PUBCOMP : Publish complete
	PUBCOMP
	// SUBSCRIBE : Subscribe
	SUBSCRIBE
	// SUBACK : Subscribe acknowledgement
	SUBACK
	// UNSUBSCRIBE : Unsubscribe
	UNSUBSCRIBE
	// UNSUBACK : Unsubscribe acknowledgement
	UNSUBACK
	// PINGREQ : Ping request
	PINGREQ
	// PINGRESP : Ping response
	PINGRESP
	// DISCONNECT : Disconnect
	DISCONNECT
)
