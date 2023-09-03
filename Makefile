PHONY: build
build:
	docker compose run --rm app goyacc -o dist/parser/parser.go parser/parser.go.y
	go build -o bin/main dist/parser/parser.go
