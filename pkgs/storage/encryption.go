package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func aesEncrypt(orig string, key string) string {
	// Convert to byte array
	origData := []byte(orig)

	k := []byte(key)

	// Group secret key
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key The length must be 16/24/32 length: %s", err.Error()))
	}
	// Gets the length of the secret key block
	blockSize := block.BlockSize()
	// Complement code
	origData = pKCS7Padding(origData, blockSize)
	// Encryption mode
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// Create array
	cryted := make([]byte, len(origData))
	// encryption
	blockMode.CryptBlocks(cryted, origData)
	//Use RawURLEncoding instead of StdEncoding
	//Do not use StdEncoding in the url parameter to cause errors
	return base64.RawURLEncoding.EncodeToString(cryted)

}

func aesDecrypt(cryted string, key string) string {
	//Use RawURLEncoding instead of StdEncoding
	//Do not use StdEncoding in the url parameter to cause errors
	crytedByte, _ := base64.RawURLEncoding.DecodeString(cryted)
	k := []byte(key)

	// Group secret key
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key The length must be 16/24/32 length: %s", err.Error()))
	}
	// Gets the length of the secret key block
	blockSize := block.BlockSize()
	// Encryption mode
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// Create array
	orig := make([]byte, len(crytedByte))
	// decrypt
	blockMode.CryptBlocks(orig, crytedByte)
	// De completion code
	orig = pKCS7UnPadding(orig)
	return string(orig)
}

//Complement
func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//De coding
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ShaGenerator(password string) string {
	fisrtHash := sha256.New()
	fisrtHash.Write([]byte(password))
	finalHash := sha256.New()
	finalHash.Write(fisrtHash.Sum(nil))
	return fmt.Sprintf("%x", finalHash.Sum(nil))
}

func generateEncKey(username string) (key string) {

	KeyHash := md5.New()
	key = string(fmt.Sprintf("%x", KeyHash.Sum([]byte(username))))

	return key[0:16]

}
