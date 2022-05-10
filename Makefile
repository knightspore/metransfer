build:
	go build -o bin/metransfer .

run:
	go run .

test: 
	go test . -v -coverprofile cover.out

cover:
	go tool cover -html=cover.out

all: test cover build
