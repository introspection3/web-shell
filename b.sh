#! /bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o amd-webshell
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o arm-webshell
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o webshell.exe
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o mac-webshell

