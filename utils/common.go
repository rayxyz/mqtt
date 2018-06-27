package utils

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// GenUUID : generate uuid
func GenUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return ""
	}
	return id.String()
}
