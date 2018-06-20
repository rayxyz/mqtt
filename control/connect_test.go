package control

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	header := new(ConnectHeader)
	payload := &ConnectPayload{
		ClientID:  "client_id看",
		WillTopic: "willtopic",
		WillMsg:   "will message",
		UserName:  "ray",
		Password:  "ray123",
	}
	payloadBytes, _ := json.Marshal(payload)
	b, err := header.Marshal(len(payloadBytes), payload)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("b => ", b)
}
