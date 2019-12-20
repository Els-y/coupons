package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

func MD5Encode(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	return encryptedData
}

func Base64Encode(data string) string {
	encodeStr := base64.StdEncoding.EncodeToString([]byte(data))
	return encodeStr
}

func Base64Decode(data string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}

func EncodeToken(username, kindStr string) string {
	return Base64Encode(username + ":" + kindStr)
}

func DecodeToken(token string) (string, string, error) {
	decodeStr, err := Base64Decode(token)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(decodeStr, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid token")
	}
	return parts[0], parts[1], nil
}
