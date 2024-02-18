# Contact Me

A simple application that collects messages from my website and sends them to my email.
It stores the messages and users in a database in case you want to do remarketing or something like that.

## Installation

```
make dockerup
make run

# TODO
# add a systemd service to restart....

```

## Uses

- Docker
- Docker Compose
- Go 1.22




## Makefile Example

```
# Makefile template

# Environment variables
DSN=username:password@tcp(localhost:3306)/dbname
MIGRATEDSN=mysql://username:password@tcp(localhost:3306)/dbname
SMTP_HOST=your_smtp_host
SMTP_PORT=your_smtp_port
SMTP_USER=your_smtp_username
SMTP_USER_PASS=your_smtp_password
SMTP_TO_ADDRESS=recipient@example.com
SMTP_FROM_ADDRESS=sender@example.com

.PHONY: docker build run migrateup migratedown deploy

docker:
	@echo "Starting Docker..."
	@docker-compose down
	@docker-compose up -d

build:
	@echo "Building..."
	@go build -o ./bin/web cmd/web/*.go

run:
	@echo "Running..."
	@go run ./cmd/web/*.go -dsn="${DSN}" -port=8080 -smtp_host="${SMTP_HOST}" -smtp_port="${SMTP_PORT}" -smtp_user="${SMTP_USER}" -smtp_user_pass="${SMTP_USER_PASS}" -smtp_to_address="${SMTP_TO_ADDRESS}" -smtp_from_address="${SMTP_FROM_ADDRESS}"

migrateup:
	@echo "Migrating up..."
	@migrate -path=./migrations -database="${MIGRATEDSN}" up

migratedown:
	@echo "Migrating down..."
	@migrate -path=./migrations -database="${MIGRATEDSN}" down

deploy:
	@echo "Deploying..."
	@$(MAKE) build
	@./bin/web/main -dsn="${DSN}" -port=8080 -smtp_host="${SMTP_HOST}" -smtp_port="${SMTP_PORT}" -smtp_user="${SMTP_USER}" -smtp_user_pass="${SMTP_USER_PASS}" -smtp_to_address="${SMTP_TO_ADDRESS}" -smtp_from_address="${SMTP_FROM_ADDRESS}"
```

## Docker Compose Example

```
version: '3.8'
services:
  mysql:
    image: mysql:latest
    restart: always
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: example_root_password
      MYSQL_DATABASE: example_database_name
      MYSQL_USER: example_user
      MYSQL_PASSWORD: example_user_password
    volumes:
      - ./data:/var/lib/mysql

```