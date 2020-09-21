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
	"os"
	"strconv"
	"time"
)

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

	if len(data)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	//fmt.Printf("%s\n", ciphertext)
	return data
}

func encrypt(key, data []byte) []byte {
	if len(data)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	//fmt.Printf("%x\n", ciphertext)
	return ciphertext
}

func main() {

	data, err := ioutil.ReadFile("1.jpg")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	file, err := os.Open("1.jpg")
	if err != nil {
		log.Printf("File open error: %s", err)
	}

	fmt.Printf("File size: %v\n", len(data))
	//sum := sha256.Sum256(data)
	//fmt.Printf("%x", sum)
	//fmt.Print(fmt.Sprintf("%x", sha256.Sum256(readBuf)))

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
	file.Close()

}
