package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// Md5 computes the MD5 hash of a string and returns it as a hexadecimal string.
func Md5(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}
