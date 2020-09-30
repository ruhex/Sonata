package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TCPCommand ...
type TCPCommand struct {
	SendFile, GetFile []byte
}

// Cmd start init byte for send to tcp stream command to server or client
var Cmd = TCPCommand{[]byte("a"), []byte("b")}

func getFile(conn net.Conn) {
	bufFileSize := make([]byte, 16)
	n, _ := conn.Read(bufFileSize)
	fmt.Printf("len %v --> file size: %s\n", n, bufFileSize)

	tmp := strings.Split(string(bufFileSize), "\n")
	buftmp := make([]byte, len(tmp[0]))
	buftmp = []byte(tmp[0])

	size, err := strconv.Atoi(fmt.Sprintf("%s", buftmp))
	if err != nil {
		fmt.Printf("Get buf size error: %s", err)
	}
	bufFile := make([]byte, size)
	conn.Read(bufFile)

	sum := sha256.Sum256(bufFile)
	fmt.Printf("SHA-256: %x\n", sum)

	file, err := os.Create(fmt.Sprintf("./files/%x", sum))
	if err != nil {
		log.Printf("Connect error: %s", err)
	}
	_, err = file.Write(bufFile)
	if err != nil {
		log.Printf("File write error: %s", err)
	}

	file.Close()
	fmt.Printf("File created\n")

	conn.Write([]byte(fmt.Sprintf("%x", sum)))
}

func sendFile(conn net.Conn) {
	bufSha256 := make([]byte, 64)
	conn.Read(bufSha256)
	data, err := ioutil.ReadFile(fmt.Sprintf("%x", bufSha256))
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("File readed")
	conn.Write([]byte(strconv.Itoa(len(data)) + "\n"))
	conn.Write(data)

}

func main() {
	fmt.Printf("Sonata server start...\n")
	fs := http.FileServer(http.Dir("./files"))
	go http.ListenAndServe(":9000", fs)
	go http.ListenAndServeTLS(":9001", "./tls/server.rsa.crt", "./tls/server.rsa.key", fs)

	ln, err := net.Listen("tcp", ":9005")
	if err != nil {
		log.Printf("%s", err)
	}
	bufCmd := make([]byte, 1)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Connect error: %s", err)
		}
		conn.Read(bufCmd)

		fmt.Printf("CMD: %x\n", string(bufCmd))

		if string(bufCmd) == "b" {
			sendFile(conn)
		}
		if string(bufCmd) == "a" {
			getFile(conn)
			// n, _ := conn.Read(bufFileSize)
			// fmt.Printf("len %v --> file size: %s\n", n, bufFileSize)

			// size, err := strconv.Atoi(fmt.Sprintf("%s", bufFileSize[0:n]))
			// if err != nil {
			// 	fmt.Printf("Get buf size error: %s", err)
			// }
			// bufFile := make([]byte, size)
			// conn.Read(bufFile)

			// sum := sha256.Sum256(bufFile)
			// fmt.Printf("SHA-256: %x\n", sum)

			// file, err := os.Create(fmt.Sprintf("./files/%x.jpg", sum))
			// if err != nil {
			// 	log.Printf("Connect error: %s", err)
			// }
			// file.Write(bufFile)

			// file.Close()
			// fmt.Printf("File created\n")

			// conn.Write([]byte(fmt.Sprintf("%x", sum)))
		}
		defer conn.Close()

	}

}
