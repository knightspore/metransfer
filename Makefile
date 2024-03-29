build:
	cd src && go build -o ./../bin/metransfer
	cp -r src/static/ bin/static/

d_build:
	docker build -t metransfer .

d_run:
	docker run -p 80:2080 -d metransfer

docker: test d_build d_run

test:
	cd src && go test . -v -coverprofile ./../cover.out
	cd src && go tool cover -html=./../cover.out


testv2:
	cd v2 && go test . -v -coverprofile ./../cover.out
	cd v2 && go tool cover -html=./../cover.out
