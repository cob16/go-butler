# go-butler
A mumble bot based on the [Gumble libary](https://github.com/layeh/gumble/)

[![Build Status](https://drone.io/github.com/njdart/go-butler/status.png)](https://drone.io/github.com/njdart/go-butler/latest)

##To run
**(Requires golang)**
- ```git clone https://github.com/njdart/go-butler.git```
- ```cd go-butler```
- ```cp ./config.json.example ./config.json```
- Edit config as necessary
- ```go run go-butler.go```

## Features
- [x] Load from config file
- [x] Load acess tokens
- [x] Get steam api status such as item servers etc [steam gauges](https://steamgaug.es/docs)
- [x] Format source connect cmds to button and give the connect to newly connected users

## Dev Features
- 'Modular' commands
- Logging

## Todo
- [ ] Functional tests agent a real mumble sever? (test bot end to end)
- [ ] check ACL's so that cmds can be made admin only
- [ ] Be able to talk back to users
