GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

init-dep:
	@dep init

dep:
	@dep ensure

status-dep:
	@dep status

update-dep:
	@dep ensure -update

test:
	@go test -v -race
	@cd ./session && go test -v -race

cover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./session && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stage/bin/service .

.PHONY: test
