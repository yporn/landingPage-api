#!/bin/bash

# ตามขั้นตอนนี้นะครับ

# 1. ทำการ sh เข้าไปที่ Container app
docker exec -it app sh

# 2. ใช้คำสั่ง Migrate Database
migrate -database 'postgres://admin:password@db:5432/sirarom_db?sslmode=disable' -source file://pkg/databases/migrations -verbose up