package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

// TCPCommand ...
type TCPCommand struct {
	SendFile, GetFile []byte
}

// Cmd start init byte for send to tcp stream command to server or client
var Cmd = TCPCommand{[]byte("a"), []byte("b")}

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
	return ciphertext
}

func main() {

	encrypting := flag.Bool("e", false, "encrypting file ")
	decrypting := flag.Bool("d", false, "decrypting file")
	local := flag.Bool("l", false, "save file local")
	fileName := flag.String("f", "", "file name")
	password := flag.String("p", "", "password")
	server := flag.String("s", "", "server ip:port")
	outFileName := flag.String("o", "config.conf", "output file name")

	flag.Parse()
	fmt.Println(*fileName)
	if *fileName == "" {
		panic("no file name!")
	}

	if *password == "" {
		panic("no password!")
	}
	tmp := *password
	data, err := ioutil.ReadFile(*fileName)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("File readed")
	//fmt.Printf("File size: %v\n", len(data))

	if *server != "" {
		conn, err := net.DialTimeout("tcp", *server, 2*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		conn.Write(Cmd.SendFile)
		conn.Write([]byte(strconv.Itoa(len(data))))
		conn.Write(data)

		if err != nil {
			log.Printf("Stream copy error: %s", err)
		}
		bufSha256 := make([]byte, 64)
		conn.Read(bufSha256)
		fmt.Printf("SHA-256: %s\n", string(bufSha256))
	}

	if *local {
		if *encrypting {
			data = encrypt([]byte(tmp), data)
			ioutil.WriteFile(fmt.Sprintf("%x", sha256.Sum256(data)), data, 0644)
		}

		if *decrypting {
			data = decrypt([]byte(tmp), data)
			ioutil.WriteFile(*outFileName, data, 0644)
		}
	}

}
