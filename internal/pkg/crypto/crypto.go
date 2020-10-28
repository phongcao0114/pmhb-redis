package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"pmhb-redis/internal/app/config"

	"golang.org/x/crypto/pbkdf2"
)

// Encrypt is for ecryption part
func Encrypt(plainText string) (string, error) {
	crypto := config.Config.Crypto
	secret := pbkdf2.Key([]byte(crypto.Password), []byte(crypto.Salt), crypto.Iteration, crypto.KeySize, sha1.New)

	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCEncrypter(cipherBlock, []byte(crypto.IV))
	plainTextWithPadding := PKCS5Padding([]byte(plainText), cipherBlock.BlockSize())

	encryptedByte := make([]byte, len(plainTextWithPadding))
	blockMode.CryptBlocks(encryptedByte, plainTextWithPadding)

	encryptedText := base64.StdEncoding.EncodeToString(encryptedByte)
	return encryptedText, nil
}

// Decrypt is for decoding data has been encrypted
func Decrypt(encryptedText string) (string, error) {
	crypto := config.Config.Crypto
	secret := pbkdf2.Key([]byte(crypto.Password), []byte(crypto.Salt), crypto.Iteration, crypto.KeySize, sha1.New)

	cipherByte, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(cipherBlock, []byte(crypto.IV))
	decryptedByte := make([]byte, len(cipherByte))
	blockMode.CryptBlocks(decryptedByte, cipherByte)

	decryptdText := PKCS5Trimming(decryptedByte)
	return string(decryptdText), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func IsEncrypted(plainText string) bool {
	_, err := Encrypt(plainText)
	if len(plainText) > 16 || err != nil {
		return true
	}
	return false
}
