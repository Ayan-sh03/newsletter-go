.PHONY: all

all:
	go run cmd/api/main.go

watch:
	nodemon