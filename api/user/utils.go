package user

import (
	"crypto/md5"
)

func GetMD5HashWithSum(text string) [16]byte {
	return md5.Sum([]byte(text))
}
