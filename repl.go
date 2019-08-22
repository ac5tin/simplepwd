package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"simplepwd/crypto"
	"strconv"
	"strings"
)

var page = 0

const printlimit int = 5

func cmdParse(cmd *[]string, data *[]standard) {
	switch (*cmd)[0] {
	case "/a":
		// add
		title := (*cmd)[1]
		username := (*cmd)[2]
		password := (*cmd)[3]
		*data = append([]standard{{Title: title, Username: username, Password: password}}, *data...)
		break
	case "/d":
		// delete
		pos := (*cmd)[1]
		posnum, err := strconv.Atoi(pos)
		if err != nil {
			break
		}
		posnum--
		*data = append((*data)[:posnum], (*data)[posnum+1:]...)
		break
	default:
		break
	}
}

func currDataPrinter(data *[]standard) {
	fmt.Printf("Page %d \n", page+1)
	for i, s := range (*data)[page*printlimit:] {
		if i == printlimit {
			break
		}
		fmt.Printf("%d. %s\n", (page*printlimit + i + 1), s.Title)
	}
}

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printStd(content *standard, r *bufio.Reader) {
	for {
		cls()
		println("=============")
		fmt.Printf("- Title    : %s\n", content.Title)
		fmt.Printf("- Username : %s\n", content.Username)
		fmt.Printf("- Password : %s\n", content.Password)
		println("=============")
		println("Enter /b to go back")
		inputraw, err := r.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		input := strings.TrimSpace(inputraw)
		if input == "/b" {
			break
		}
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

	reader := bufio.NewReader(os.Stdin)
	for {
		cls()
		println("simplepwd 📝")
		println("============")
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
			cls()
			break
		case "/s":
			// save
			println("-- save --")
			d, err := json.Marshal(arrdata)
			if err != nil {
				panic(err.Error())
			}
			crypto.EncryptFile(*filename, &d, *password)
			break
		case "/n":
			// next page
			if (len(arrdata) / printlimit) >= (page + 1) {
				page++
			}
			break
		case "/p":
			// previous
			if page > 0 {
				page--
			}
			break
		default:
			// parse
			cmd := strings.Split(input, " ")
			switch len(cmd) {
			case 0:
				break
			case 1:
				// probably number
				i, err := strconv.Atoi(cmd[0])
				if err != nil {
					break
				}
				printStd(&arrdata[i-1], reader)
				break
			default:
				cmdParse(&cmd, &arrdata)
				break
			}
			continue
		}
	}
}
