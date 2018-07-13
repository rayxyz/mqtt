package utils

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println("uuid => ", GenUUID())
}

func TestRemoveSliceEles(t *testing.T) {
	ages := []int{3, 45, 23, 423, 45, 2, 3}
	fmt.Println(ages[1:3])
	for i := len(ages) - 1; i >= 0; i-- {
		if ages[i] > 20 {
			ages = append(ages[:i], ages[i+1:]...)
		}
	}
	fmt.Println(ages)
	agesx := []int{3, 3}
	fmt.Println(agesx[2:])
}
