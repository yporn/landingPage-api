#!/bin/bash


# /go/bin/
migrate -database 'postgres://admin:password@db:5432/sirarom_db?sslmode=disable' -source file://pkg/databases/migrations -verbose up

# Start the main application
./main
