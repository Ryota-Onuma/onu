PHONY: build
build:
	docker compose run --rm app goyacc -o core/parser.go core/parser.y
	go build -o bin/onu ./main.go
