package main

import (
	"flag"
	"fmt"
	"simplepwd/crypto"
	"simplepwd/useful"
)

var (
	filename *string
	password *string
	dec      *bool
)

func main() {
	filename = flag.String("f", "", "File path of passwords")
	password = flag.String("p", "", "Master password of file")
	dec = flag.Bool("d", false, "Decrypt file and print content")
	flag.Parse()

	if useful.FileExist(*filename) {
		// file exists, decrypt file content
		data := crypto.DecryptFile(*filename, *password)
		if *dec {
			fmt.Println(string(*data))
			return
		}
		repl(data)
	} else {
		// file doesnt exist, create file
		useful.CreateFile(*filename)
	}
}
