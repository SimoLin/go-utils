package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strings"
)

func MD5Encode(S1 string) string {
	ret := md5.Sum([]byte(S1))
	return hex.EncodeToString(ret[:])
}

func MD5EncodeByte(B1 []byte) string {
	ret := md5.Sum(B1)
	return hex.EncodeToString(ret[:])
}

func MD5EncodeFile(file_path string) (r string) {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return
	}
	return MD5EncodeByte(data)
}

func SHA1Encode(S1 string) string {
	o := sha1.New()
	o.Write([]byte(S1))
	return hex.EncodeToString(o.Sum(nil))
}

func SHA256Encode(S1 string) string {
	o := sha256.New()
	o.Write([]byte(S1))
	return hex.EncodeToString(o.Sum(nil))
}

func HMACSHA512Encode(s string, key string) string {
	o := hmac.New(sha512.New, []byte(key))
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

func Base64Encode(S1 string) string {
	return base64.StdEncoding.EncodeToString([]byte(S1))
}

func Base64Decode(S1 string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(S1)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

// 尝试 Base64 解码，失败时返回 原字符串 而不是 空
func TryBase64Decode(S1 string) (r string) {
	if i := len(S1) % 4; i != 0 {
		S1 += strings.Repeat("=", 4-i)
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(S1)
	if err != nil {
		return S1
	}
	return string(decodeBytes)
}

// 尝试 Base64 RawURL 解码，失败时返回 原字符串 而不是 空
func TryBase64DecodeRawURL(S1 string) (r string) {
	if i := len(S1) % 4; i != 0 {
		S1 += strings.Repeat("=", 4-i)
	}
	decodeBytes, err := base64.RawURLEncoding.DecodeString(S1)
	if err != nil {
		return S1
	}
	return string(decodeBytes)
}
