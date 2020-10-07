package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

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
		if err != nil {
			fmt.Printf("Connect error: %s", err)
		}
		bufCmd := make([]byte, 1)
		conn.Read(bufCmd)
		fmt.Printf("%s", bufCmd)
		conn.Write([]byte("ok"))
		conn.Close()
	}
}
