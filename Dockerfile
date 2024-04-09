FROM --platform=linux/amd64 golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./


RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN apk update && apk add --no-cache curl
# Install migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

RUN go mod tidy

COPY . .

RUN go build -o main .

# RUN chmod +x entrypoint.sh

FROM --platform=linux/amd64 alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/pkg/databases/migrations /app/pkg/databases/migrations


# RUN migrate -database 'postgres://admin:password@db:5432/sirarom_db?sslmode=disable' -source file://pkg/databases/migrations -verbose up
EXPOSE 3000

# CMD ["/app/main"]
# CMD ["./app/main"]
ENTRYPOINT [ "/app/main" ]
