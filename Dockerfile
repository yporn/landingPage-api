FROM golang:latest AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y postgresql

COPY go.mod go.sum ./

RUN go mod download

RUN go get github.com/lib/pq@latest
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go mod tidy

COPY . .

RUN go build -o main .

RUN chmod +x entrypoint.sh

EXPOSE 3000

CMD ["./entrypoint.sh"]
