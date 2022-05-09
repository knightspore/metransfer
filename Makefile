build:
	go build -o bin/metransfer .

run:
	go run .

test: 
	go test . -v -coverprofile cover.out
	go tool cover -html=cover.out

all: test build
