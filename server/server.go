package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

// TCPCommand ...
type TCPCommand struct {
	SendFile, GetFile []byte
}

// Cmd start init byte for send to tcp stream command to server or client
var Cmd = TCPCommand{[]byte("a"), []byte("b")}

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
	bufFileSize := make([]byte, 16)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Connect error: %s", err)
		}
		conn.Read(bufCmd)

		fmt.Printf("CMD: %x\n", string(bufCmd))

		if string(bufCmd) == "a" {

			n, _ := conn.Read(bufFileSize)
			fmt.Printf("len %v --> file size: %s\n", n, bufFileSize)

			size, err := strconv.Atoi(fmt.Sprintf("%s", bufFileSize[0:n]))
			if err != nil {
				fmt.Printf("Get buf size error: %s", err)
			}
			bufFile := make([]byte, size)
			conn.Read(bufFile)

			sum := sha256.Sum256(bufFile)
			fmt.Printf("SHA-256: %x\n", sum)

			file, err := os.Create(fmt.Sprintf("./files/%x.jpg", sum))
			if err != nil {
				log.Printf("Connect error: %s", err)
			}
			file.Write(bufFile)

			file.Close()
			fmt.Printf("File created\n")

			/* data, err := ioutil.ReadFile("test.jpg")
			if err != nil {
				fmt.Printf("File reading error: %s", err)
				return
			} */
			conn.Write([]byte(fmt.Sprintf("%x", sum)))
			//fmt.Fprintf(conn, "1\n")
		}
		defer conn.Close()

	}

}
