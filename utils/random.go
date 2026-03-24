package utils

import (
	"github.com/google/uuid"
	"strings"
)

// RandomString generates a random string using a UUID v4 with hyphens removed.
func RandomString() string {
	newUUID, _ := uuid.NewRandom()
	return strings.Replace(newUUID.String(), "-", "", -1)
}
