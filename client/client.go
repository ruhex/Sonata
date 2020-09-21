package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

// ----------------------------- AES CFB Decrypt ----------------------------- //
func decrypt(key, data []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(data) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	//fmt.Printf("%s", data)
	return data
}

// ----------------------------- AES CFB Encrypt ----------------------------- //
func encrypt(key, data []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	//fmt.Printf("%x\n", ciphertext)
	return ciphertext
}

func main() {

	data, err := ioutil.ReadFile("404.jpg")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	fmt.Printf("File size: %v\n", len(data))

	conn, err := net.DialTimeout("tcp", "127.0.0.1:9005", 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("a"))
	conn.Write([]byte(strconv.Itoa(len(data))))
	conn.Write(data)

	if err != nil {
		log.Printf("Stream copy error: %s", err)
	}
	bufSha256 := make([]byte, 64)
	conn.Read(bufSha256)
	fmt.Printf("SHA-256: %s\n", string(bufSha256))

}
