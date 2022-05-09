build:
	go build -o bin/metransfer .

run:
	go run .

test: 
	go test . -v -cover

all: test build
