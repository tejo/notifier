FROM golang:latest as builder

WORKDIR $GOPATH/src/

COPY . .

RUN make build

RUN ["./notify", "-h"]

# docker build -t notifier .
