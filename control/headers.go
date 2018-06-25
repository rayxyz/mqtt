package control

// PubAckHeader : Publish acknowledgement header
type PubAckHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// PubRecHeader : Publish receive header
type PubRecHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// PubRelHeader : Publish release header
type PubRelHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// PubCompHeader : Publish complete header
type PubCompHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// SubHeader : Subscribe header
type SubHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// SubPayload : Subscribe payload
type SubPayload struct {
	TopicFilter []byte
	ReqQoS      int
}

// SubAckHeader : Subscribe acknowledgement header
type SubAckHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// SubAckPayload : Subscribe acknowledgement payload
type SubAckPayload struct {
	ReturnCode int
}

// UnsubHeader : Unsubscribe header
type UnsubHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// UnsubPayload : Unsubscribe payload
type UnsubPayload []byte

// UnsubAckHeader : Unsubscribe acknowledgement header
type UnsubAckHeader struct {
	PackType  int
	RemainLen int
	PackID    int
}

// PingReqHeader : Ping request header
type PingReqHeader struct {
	PackType  int
	ReaminLen int
	PackID    int
}

// PingRespHeader : Ping response header
type PingRespHeader struct {
	PackType  int
	RemainLen int
}

// DisconnectHeader : Disconnect connection header
type DisconnectHeader struct {
	PackType  int
	RemainLen int
}
