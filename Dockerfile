FROM golang:1.15 as build-env

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o web
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./cmd/migrate 


# final stage
FROM alpine

RUN apk add ca-certificates

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=build-env /go/src/app/web .
COPY --from=build-env /go/src/app/migrate .

EXPOSE 80

CMD ["./web"]
