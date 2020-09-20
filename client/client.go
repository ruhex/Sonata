package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

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

	conn, _ := net.Dial("tcp", "127.0.0.1:9002")

	conn.Write([]byte("file_send\n"))
	io.Copy(conn, file)

	conn.Close()

}
