package control

import (
	"fmt"
	"testing"
)

func TestTypes(t *testing.T) {
	fmt.Println("CONNECT => ", CONNECT)
	fmt.Println("CONNACK => ", CONNACK)
	fmt.Println("PUBLISH => ", PUBLISH)
	fmt.Println("PUBACK => ", PUBACK)
	fmt.Println("PUBREC => ", PUBREC)
	fmt.Println("PUBREL => ", PUBREL)
	fmt.Println("PUBCOMP => ", PUBCOMP)
	fmt.Println("SUBSCRIBE => ", SUBSCRIBE)
	fmt.Println("SUBACK => ", SUBACK)
	fmt.Println("UNSUBSCRIBE => ", UNSUBSCRIBE)
	fmt.Println("UNSUBACK => ", UNSUBACK)
	fmt.Println("PINGREQ => ", PINGREQ)
	fmt.Println("PINGRESP => ", PINGRESP)
	fmt.Println("DISCONNECT => ", DISCONNECT)
}
