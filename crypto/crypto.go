package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func AESEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted
}

func PKCS7Padding(originByte []byte, blockSize int) []byte {
	padding := blockSize - len(originByte)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(originByte, padText...)
}

// 解析 string 为 rsa.PublicKey 类型
func RSAReadPublicKey(pub_key_string string) (pub_key *rsa.PublicKey, err error) {
	pemBlock, _ := pem.Decode([]byte(pub_key_string))
	var pkixPublicKey any
	if pemBlock.Type == "RSA PUBLIC KEY" {
		// PKCS1类型以 -----BEGIN RSA PUBLIC KEY----- 开头
		pkixPublicKey, err = x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	} else if pemBlock.Type == "PUBLIC KEY" {
		// PKIX类型以 -----BEGIN PUBLIC KEY----- 开头
		pkixPublicKey, err = x509.ParsePKIXPublicKey(pemBlock.Bytes)
	}
	if err != nil {
		return nil, err
	}
	publicKey := pkixPublicKey.(*rsa.PublicKey)
	return publicKey, nil
}

// 解析 string 为 rsa.PrivateKey 类型
func RSAReadPrivateKey(priv_key_string string) (priv_key *rsa.PrivateKey, err error) {
	pemBlock, _ := pem.Decode([]byte(priv_key_string))
	// PKCS1类型以 -----BEGIN RSA PUBLIC KEY----- 开头
	priv_key, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return
}

func RSAEncrypt(plain_text string, pub_key *rsa.PublicKey) (encrypt_text string, err error) {
	encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, pub_key, []byte(plain_text))
	if err != nil {
		return "", err
	}
	encrypt_text = base64.StdEncoding.EncodeToString(encryptPKCS1v15)
	return
}

func RSADecrypt(encrypt_text string, priv_key *rsa.PrivateKey) (plain_text string, err error) {
	decryptPKCS1v15, err := rsa.DecryptPKCS1v15(rand.Reader, priv_key, []byte(encrypt_text))
	if err != nil {
		return "", err
	}
	plain_text = string(decryptPKCS1v15)
	return
}
