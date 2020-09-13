FROM golang:1.15

WORKDIR $GOPATH/src/app

COPY . .

RUN go mod tidy

# RUN go get github.com/matryer/moq
