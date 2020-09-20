package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {

	data, err := ioutil.ReadFile("1.jpg")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Printf("File size: %v\n", len(data))
	//sum := sha256.Sum256(data)
	//fmt.Printf("%x", sum)
	//fmt.Print(fmt.Sprintf("%x", sha256.Sum256(readBuf)))

	conn, _ := net.Dial("tcp", "127.0.0.1:9002")

	conn.Write([]byte("cmd test"))

	conn.Close()

}
