package common

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5HashWithSum(text string) string {
	byteSum := md5.Sum([]byte(text))
	return hex.EncodeToString(byteSum[:])
}
