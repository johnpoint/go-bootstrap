package utils

import (
	"OwlBackend/config"
	"crypto/sha256"
	"fmt"
)

func EncodePassword(password, username string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(config.Config.Salt+password+username)))
}

func EqualPassword(password, standardPwd, username string) bool {
	return standardPwd == EncodePassword(password, username)
}
