FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8500

CMD ["./main"]
