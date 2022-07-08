test: 
	cd src && go test . -v -coverprofile ./../cover.out

cover:
	cd src && go tool cover -html=./../cover.out

build:
	cd src && go build -o ./../bin/metransfer
	cp -r src/static/ bin/static/

dev:
	cd src && go run .

compile:
	env GOOS=windows GOARCH=amd64 cd src && go build -o ./../bin/metransfer.exe
	env GOOS=darwin GOARCH=amd64 cd src && go build -o ./../bin/metransfer.darwin.amd64
	env GOOS=linux GOARCH=amd64 cd src && go build -o ./../bin/metransfer

run:
	go run .

dockerbuild:
	docker build -t metransfer .

dockerrun:
	docker run -p 8080:2080 -d metransfer

all: test cover build compile
