package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"simplepwd/crypto"
	"strconv"
	"strings"
)

var page = 0
var searchTxt = ""

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
	case "/u":
		// update
		pos := (*cmd)[1]
		field := (*cmd)[2]
		value := (*cmd)[3]
		posnum, err := strconv.Atoi(pos)
		if err != nil {
			break
		}
		posnum--

		updated := (*data)[posnum]
		switch field {
		case "Title", "title":
			updated.Title = value
			break
		case "Username", "username":
			updated.Username = value
			break
		case "Password", "password":
			updated.Password = value
			break
		default:
			break
		}
		(*data)[posnum] = updated
		break
	case "/f":
		// find/search
		searchTxt = (*cmd)[1]
		break
	default:
		break
	}
}

func filter(data *[]standard) []standard {
	var retme []standard
	for _, v := range *data {
		// search text
		if strings.Contains(v.Title, searchTxt) {
			retme = append(retme, v)
		}
	}
	return retme
}

func currDataPrinter(data *[]standard) {
	filtered := filter(data)
	fmt.Printf("Page %d of %d\n", page+1, len(filtered)/printlimit+1)
	for i, s := range (filtered)[page*printlimit:] {
		if i == printlimit {
			break
		}
		fmt.Printf("%d. %s\n", (page*printlimit + i + 1), s.Title)
	}
}

func cls() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		break
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		break
	}

}

func printFilter() {
	if searchTxt != "" {
		fmt.Printf("Searching : %s (Enter /f to reset)\n", searchTxt)
	}
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
		printFilter()
		currDataPrinter(&arrdata)
		print("> ")
		inputraw, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		input := strings.TrimSpace(inputraw)
		switch input {
		// exact matches (command wiht no arguments)
		case ":q", "/q", "bye", "quit", "exit":
			println("Bye ~")
			return
		case "clear":
			cls()
			break
		case "/s":
			// save/write
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
		case "/f":
			// find/search (CLEAR)
			searchTxt = ""
			break
		default:
			// parse
			// first check there are spaces
			// if spaces exist, probably means user wants to enter a command (not select a number)
			cmd := strings.Split(input, " ")
			switch len(cmd) {
			case 0:
				// user didnt enter anything,
				break
			case 1:
				// user only entered 1 argument, no spaces
				// probably number
				i, err := strconv.Atoi(cmd[0])
				if err != nil {
					// no, user didnt enter a number
					break
				}
				// user entered a number, print info
				filteredData := filter(&arrdata)
				printStd(&filteredData[i-1], reader)
				break
			default:
				// user entered more than 1 argument, probably a command (with arguments)
				cmdParse(&cmd, &arrdata)
				break
			}
			continue
		}
	}
}
