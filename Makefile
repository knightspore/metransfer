build:
	go build -o bin/metransfer .

compile:
	env GOOS=windows GOARCH=amd64 go build -o bin/metransfer.exe
	env GOOS=darwin GOARCH=amd64 go build -o bin/metransfer.darwin.amd64
	env GOOS=linux GOARCH=amd64 go build -o bin/metransfer .

run:
	go run .

test: 
	go test . -v -coverprofile cover.out

cover:
	go tool cover -html=cover.out

all: test cover build
