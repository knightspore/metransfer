FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN cd /build && git clone https://github.com/knightspore/metransfer.git
RUN cd /build/metransfer && make build

EXPOSE 2080

CMD [ "/build/metransfer/bin/metransfer" ]

