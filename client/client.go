package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"
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

	conn, err := net.DialTimeout("tcp", "127.0.0.1:9005", 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("a"))
	conn.Write([]byte(strconv.Itoa(len(data))))
	conn.Write(data)

	//conn.SetReadDeadline(time.Now().Add(time.Second))
	//count, err := io.Copy(conn, file)
	if err != nil {
		log.Printf("Stream copy error: %s", err)
	}

	//fmt.Printf("count write to stream: %v\n", count)
	bufSha256 := make([]byte, 64)
	conn.Read(bufSha256)
	fmt.Printf("CMD: %s\n", string(bufSha256))
	file.Close()

}
