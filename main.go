package main

import (
	"flag"
	"fmt"
	"log"
	"simplepwd/crypto"
	"simplepwd/useful"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	filename *string
	password *string
	dec      *bool
	chgpw    *bool
)

func main() {
	filename = flag.String("f", "", "File path of passwords")
	password = flag.String("p", "", "Master password of file")
	dec = flag.Bool("d", false, "Decrypt file and print content")
	chgpw = flag.Bool("c", false, "Change Passowrd")
	flag.Parse()

	if useful.FileExist(*filename) {
		// file exists, decrypt file content
		if *password == "" {
			fmt.Println("Please Enter Password")
			bytePw, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatal(err.Error())
			}
			*password = string(bytePw)
		}

		data := crypto.DecryptFile(*filename, *password)
		if *dec {
			fmt.Println(string(*data))
			return
		}
		if *chgpw {
			fmt.Println("Please Enter New Password")
			bytePw, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatal(err.Error())
			}
			crypto.EncryptFile(*filename, data, string(bytePw))
		}

		repl(data)
	} else {
		// file doesnt exist, create file
		useful.CreateFile(*filename)
	}
}
