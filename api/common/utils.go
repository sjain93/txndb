package common

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

func GetMD5HashWithSum(text string) string {
	byteSum := md5.Sum([]byte(text))
	return hex.EncodeToString(byteSum[:])
}

func StringToCents(amount string) (int64, error) {
	a := strings.TrimSpace(amount)
	a = strings.Replace(a, ",", "", -1)

	float, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	return int64(float * 100), nil
}

func UnmarshalTime(csv, format string) (time.Time, error) {
	return time.Parse(format, csv)
}
