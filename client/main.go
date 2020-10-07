package main

import (
	"flag"
	"io/ioutil"
	"net"
	"time"

	"../pkg/crypt"
)

var (
	//BuildInfo string = ``

	argv struct {
		//port    uint
		server  string
		local   bool
		decrypt bool
		encrypt bool
		file    string
		outname string
		passwd  string
		help    bool
	}
)

func init() {
	flag.StringVar(&argv.server, `s`, `localhost:9005`, `remote conf server`)
	flag.BoolVar(&argv.help, `h`, false, `show this help`)
	flag.BoolVar(&argv.local, `l`, false, `enable local save file`)
	flag.BoolVar(&argv.decrypt, `d`, false, `decrypt file`)
	flag.BoolVar(&argv.encrypt, `e`, false, `necrypt file`)
	flag.StringVar(&argv.file, `f`, ``, `open file name`)
	flag.StringVar(&argv.outname, `o`, `config.conf`, `out file name`)
	flag.StringVar(&argv.passwd, `p`, ``, `password for crypt`)
	flag.Parse()
}

func main() {
	if argv.help {
		flag.Usage()
		return
	}

	//sendData([]byte("a"), 10)

	file, err := fileRead(argv.file)
	if err != nil {
		println(err)
	}

	if argv.encrypt {
		file = crypt.Encrypt(argv.passwd, file)
	}

	if argv.decrypt {
		file = crypt.Decrypt(argv.passwd, file)
	}

	if savefile(argv.outname, file) != nil {
		println(err)
	}

}

func fileRead(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func savefile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func sendData(cmd []byte, buffLen int64) ([]byte, error) {
	buff := make([]byte, buffLen)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9005", 2*time.Second)
	defer conn.Close()

	if err != nil {
		return nil, err

	}

	conn.Write(cmd)
	conn.Read(buff)
	return buff, nil

}
