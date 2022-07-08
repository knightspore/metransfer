FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN cd /build && git clone https://github.com/knightspore/metransfer.git
RUN cd /build/metransfer && go build -o /metransfer
COPY index.html /

EXPOSE 2080

CMD [ "/metransfer" ]

