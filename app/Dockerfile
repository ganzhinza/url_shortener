FROM golang:1.24 AS builder

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd/url_shortener
RUN GOOS=linux go build -o app

EXPOSE 8080
CMD ["./app"]