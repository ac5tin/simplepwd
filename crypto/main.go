package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func enc(data *[]byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, *data, nil)
	return ciphertext
}

func dec(data *[]byte, passphrase string) *[]byte {
	if len(*data) == 0 {
		retme := []byte("")
		return &retme
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := (*data)[:nonceSize], (*data)[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return &plaintext
}

// EncryptFile encrypts a file with a password
func EncryptFile(filename string, data *[]byte, passphrase string) {
	f, _ := os.OpenFile(filename, os.O_RDWR, 0644)
	defer f.Close()
	f.Write(enc(data, passphrase))
}

// DecryptFile decrypts file with password
func DecryptFile(filename string, passphrase string) *[]byte {
	data, _ := ioutil.ReadFile(filename)
	return dec(&data, passphrase)
}
