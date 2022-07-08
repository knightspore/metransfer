FROM golang:latest

RUN mkdir /build

RUN cd /build && git clone https://github.com/knightspore/metransfer.git
RUN cd /build/metransfer && make build

WORKDIR /build/metransfer/bin

EXPOSE 2080

CMD [ "/build/metransfer/bin/metransfer" ]

