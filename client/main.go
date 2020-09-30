package main

import (
	"flag"
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
	flag.BoolVar(&argv.decrypt, `d`, false, `encrypt file`)
	flag.BoolVar(&argv.encrypt, `e`, false, `decrypt file`)
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
}
