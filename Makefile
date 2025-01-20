.PHONY: setup gen test

setup:
	go install go.uber.org/mock/mockgen@latest

gen:
	go generate ./...

test: setup gen
	go test ./... -cover -race
