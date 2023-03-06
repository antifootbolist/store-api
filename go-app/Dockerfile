FROM golang:1.20-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go-app/go.mod go-app/go.sum .
RUN go mod download

COPY go-app/main.go .
RUN go build -o go-app main.go

EXPOSE 8080

CMD ["./go-app"]
