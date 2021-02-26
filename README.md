# goc2
c2 client/server/paylod

![](https://github.com/grines/goc2/blob/main/goc2.gif)

# GoC2 - MacOS Post Exploitation C2 Framework

Custom C2 for bypassing EDR and ease of use.

## Status
- This is still an active work in progress (Not ready for production use.. I made it in a weekend.. has bugs.)

## Features
- [x] Terraform deployment
- [X] Command History
- [X] Remote Command Completion (yes this works!)
- [X] JXA execution (cocoa api)
- [X] Clipboard (cocoa api)
- [X] cat / curl (cocoa api)
- [ ] add Doom persistence list
- [ ] Add Slack integration
- [ ] Add ++ persistence
- [ ] Add + privesc
- [ ] Encrytpion
- [ ] variable callback timeout
- [ ] Authentication
- [ ] Custom JXA paylaods storage

## Prereqs
- install mongodb on c2 server ** sudo apt install mongodb * required
 
## Getting Started (C2 Server)
- go get github.com/goc2
- sudo apt install mongodb || brew install mongodb
- ./goc2 --web

## CLI
- ./goc2 --cli --c2 http://c2.server 

## Payloads
- grab a goc2-agent macos payload
- edit c2 ip before compiling
- ./agent
