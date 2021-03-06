# SimplePassword 🔑
## very minimal password manager

<a href="https://asciinema.org/a/zC84FUt0hazWae8QJMnJz9Xjw" target="_blank"><img src="https://asciinema.org/a/zC84FUt0hazWae8QJMnJz9Xjw.svg" /></a>

## Usage
### flags
```
-f  filename
-p  password
-d  decrypt file and print out content (optional, default=false)
```

### To Start
```sh
simplepwd -f abc -p "abc123"

# all in one line
simplepwd -f abc -p "abc123"
```

### add record
```
/a TITLE USERNAME PASSWORD
```
### delete record
```
/d 3
```
### update record
```
/u INDEX FIELD VALUE
FIELD = title / username / password
```
#### e.g.
```
/u 23 username john
/u 23 password abcdef1
```

### save
```
/s
```

### find / search record
```
/f foobar
```
#### To reset search text
```
/f
```


### navigate pages
```
/n  next page
/p  previous page
```


### show info
```
2   number/index of the record
```

### Copy to clipboard (while showing info)
```
pw  # copies password
user # copies username
```

### Update / Change File Encryption password
```
simplepwd -f pass -c
```


### Print file content and read it in less
```
simplepwd -f abc -p "abc123" -d|jq -C .|less -R
```

### Run
```
go run *.go -f abc -p "abc123"
```
### Build + Run
```
go build -o bin/simplepwd
./bin/simplepwd -f abc -p "abc123"
```

## Progress
- [x] Encrypt/Decrypt
- [x] Basic REPL
- [x] Add record
- [x] Save
- [x] Output to stdout
- [x] Show password in REPL
- [x] Remove record
- [x] Edit record
- [x] Search record
- [x] Copy to clipboard
- [ ] Research stronger Encrypt/Decrypt
- [ ] Error handling
- [ ] More/Custom Fields
