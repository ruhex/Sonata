package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

func getSize(data []byte) int {
	for index := range data {
		if data[index] == '\x00' {
			return index
		}
	}
	return 0
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

	for {
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Printf("Connect error: %s", err)
		}
		bufCmd := make([]byte, 1)
		buffSize := make([]byte, 16)
		conn.Read(bufCmd)
		if string(bufCmd) == "a" {
			println("command: a")
			conn.Read(buffSize)
			size, err := strconv.Atoi(fmt.Sprintf("%s", buffSize[0:getSize(buffSize)]))
			if err != nil {
				fmt.Printf("error: %v", err)
			}
			fmt.Printf("size: %v\n", size)
			buffFile := make([]byte, size)

			n, err := conn.Read(buffFile)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
			fmt.Printf("len1: %v", n)

			fmt.Printf("buffer: %v\n", len(buffFile))
			ioutil.WriteFile(fmt.Sprintf("./files/%x", sha256.Sum256(buffFile)), buffFile, 0644)
		}
		conn.Write([]byte("ok"))

	}
}
