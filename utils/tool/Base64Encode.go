package utils

import (
	"encoding/base64"
)

func Base64Decode(str string) []byte {
	var b []byte
	var err error
	x := len(str) * 3 % 4
	switch {
	case x == 2:
		str += "=="
	case x == 1:
		str += "="
	}
	if b, err = base64.StdEncoding.DecodeString(str); err != nil {
		return b
	}

	return b
}
