package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Printf("Sonata server start...\n")
	//fs := http.FileServer(http.Dir("./files"))
	//go http.ListenAndServe(":9000", fs)
	//go http.ListenAndServeTLS(":9001", "server.rsa.crt", "server.rsa.key", fs)

	ln, err := net.Listen("tcp", ":9005")
	if err != nil {
		log.Printf("%s", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Connect error: %s", err)
		}

		cmd, _ := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Reader stream error: %s", err)
		}
		fmt.Printf("CMD: %s\n", string(cmd))

		if string(cmd) == "file_send\n" {
			file, err := os.Create("test.jpg")
			if err != nil {
				log.Printf("Connect error: %s", err)
			}
			io.Copy(file, conn)
			file.Close()
			fmt.Printf("File created\n")
			data, err := ioutil.ReadFile("test.jpg")
			if err != nil {
				fmt.Printf("File reading error: %s", err)
				return
			}
			sum := sha256.Sum256(data)
			fmt.Printf("SHA-256: %x\n", sum)
			//conn.Write([]byte(fmt.Sprintf("%x\n", sum)))
		}
		defer conn.Close()

	}

}
