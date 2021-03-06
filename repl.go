package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"simplepwd/crypto"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
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
	case "/d":
		// delete
		pos := (*cmd)[1]
		posnum, err := strconv.Atoi(pos)
		if err != nil {
			break
		}

		*data = append((*data)[:posnum], (*data)[posnum+1:]...)
	case "/u":
		// update
		pos := (*cmd)[1]
		field := (*cmd)[2]
		value := (*cmd)[3]
		posnum, err := strconv.Atoi(pos)
		if err != nil {
			break
		}

		updated := (*data)[posnum]
		switch field {
		case "Title", "title":
			updated.Title = value
		case "Username", "username":
			updated.Username = value
		case "Password", "password":
			updated.Password = value
		default:
			break
		}
		(*data)[posnum] = updated
	case "/f":
		// find/search
		page = 0
		searchTxt = (*cmd)[1]
	default:
		break
	}
}

func filter(data *[]standard) map[uint32]standard {
	retme := make(map[uint32]standard)
	for i, v := range *data {
		// search text
		if strings.Contains(v.Title, searchTxt) {
			retme[uint32(i)] = v
		}
	}
	return retme
}

func currDataPrinter(data *[]standard) {
	filtered := filter(data)
	fmt.Printf("Page %d of %d\n", page+1, len(filtered)/printlimit+1)

	// sort key
	keys := []uint32{}
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(x, y int) bool { return keys[x] < keys[y] })

	i := 0
	for _, k := range keys {
		if i < page*printlimit {
			i++
			continue
		}
		if i >= page*printlimit+printlimit {
			break
		}
		fmt.Printf("%d. %s\n", k, filtered[k].Title)
		i++
	}
}

func cls() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func printFilter() {
	if searchTxt != "" {
		fmt.Printf("Searching : %s (Enter /f to reset)\n", searchTxt)
	}
}

func pwHide(pw string) string {
	retme := ""
	for i := 0; i < len(pw); i++ {
		retme += "*"
	}
	return retme
}

func printStd(content *standard, r *bufio.Reader) {
loop:
	for {
		cls()
		println("=============")
		fmt.Printf("- Title    : %s\n", content.Title)
		fmt.Printf("- Username : %s\n", content.Username)
		fmt.Printf("- Password : %s\n", pwHide(content.Password))
		println("=============")
		println("Enter /b to go back")
		inputraw, err := r.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		input := strings.TrimSpace(inputraw)
		switch input {
		case "user", "username":
			clipboard.WriteAll(content.Username)
		case "pw", "password":
			clipboard.WriteAll(content.Password)
		case "/b":
			break loop
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
		case "/s":
			// save/write
			println("-- save --")
			d, err := json.Marshal(arrdata)
			if err != nil {
				panic(err.Error())
			}
			crypto.EncryptFile(*filename, &d, *password)
		case "/n":
			// next page
			if (len(arrdata) / printlimit) >= (page + 1) {
				page++
			}
		case "/p":
			// previous
			if page > 0 {
				page--
			}
		case "/f":
			// find/search (CLEAR)
			searchTxt = ""
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
				t := filteredData[uint32(i)]
				printStd(&t, reader)
			default:
				// user entered more than 1 argument, probably a command (with arguments)
				cmdParse(&cmd, &arrdata)
			}
			continue
		}
	}
}
