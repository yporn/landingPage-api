FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# RUN go get github.com/lib/pq@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go mod tidy

COPY . .

RUN go build -o main .

RUN chmod +x entrypoint.sh

# EXPOSE 3000

CMD ["./entrypoint.sh"]
