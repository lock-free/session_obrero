GO111MODULE := on
export GO111MODULE

init:
	@go mod init

clean:
	@go mod tidy

update:
	@go get -u

run:
	@go run main.go

test:
	@go test -v -race
	@cd ./session && go test -v -race

cover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./session && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stage/bin/service .

.PHONY: test
