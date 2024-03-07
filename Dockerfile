FROM golang:1.20-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/app

EXPOSE 3000

ENTRYPOINT [ "/bin/app", "/bin/.env.prod" ]