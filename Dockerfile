FROM golang:1.22

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN git config --global --add safe.directory /workdir
RUN go mod tidy