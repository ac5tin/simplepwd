# SimplePassword 🔑
## very minimal password manager

## Usage
### flags
```
-f  filename
-p  password
-d  decrypt file and print out content (optional, default=false)
```

### To Start
```
simplepwd -f abc -p "abc123"
```

### add record
```
/a TITLE USERNAME PASSWORD
```

### Print file content and read it in less
```
simplepwd -f abc -p "abc123" -d|jq -C .|less -R
```


## Progress
- [x] Encrypt/Decrypt
- [x] Basic REPL
- [x] Add record
- [x] Save
- [x] Output to stdout
- [ ] Show passwor in REPL
- [ ] Remove record
- [ ] Edit record
- [ ] Research stronger Encrypt/Decrypt
- [ ] Error handling
