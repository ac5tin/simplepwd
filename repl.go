package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"simplepwd/crypto"
	"strings"
)

func cmdParse(cmd *[]string, data *[]standard) {
	switch (*cmd)[0] {
	case "/a":
		// add
		title := (*cmd)[1]
		username := (*cmd)[2]
		password := (*cmd)[3]
		*data = append([]standard{{Title: title, Username: username, Password: password}}, *data...)
		break
	default:
		break
	}
}

func currDataPrinter(data *[]standard) {
	printlimit := 5
	for i, s := range *data {
		if i == printlimit {
			break
		}
		fmt.Printf("%d. %s\n", i, s.Title)
	}
}

func repl(data *[]byte) {
	arrdata := []standard{}
	if len(*data) > 0 {
		err := json.Unmarshal(*data, &arrdata)
		if err != nil {
			panic(err.Error())
		}
	}

	println("simplepwd 📝")
	println("============")
	reader := bufio.NewReader(os.Stdin)
	for {
		currDataPrinter(&arrdata)
		print("> ")
		inputraw, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		input := strings.TrimSpace(inputraw)
		switch input {
		case ":q", "/q", "bye", "quit", "exit":
			println("Bye ~")
			return
		case "clear":
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			break
		case "/s":
			// save
			d, err := json.Marshal(arrdata)
			if err != nil {
				panic(err.Error())
			}
			crypto.EncryptFile(*filename, &d, *password)
			println("-- saved --")
			break
		default:
			// parse
			cmd := strings.Split(input, " ")
			if len(cmd) < 2 {
				continue
			}
			cmdParse(&cmd, &arrdata)
			continue
		}
	}
}
