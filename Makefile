# Makefile

SERVER_TCP_ADDRESS ?= localhost:4200
SERVER_UDP_ADDRESS ?= localhost:4201
PUBLIC_TCP_ADDRESS ?= localhost:4200
PUBLIC_UDP_ADDRESS ?= localhost:4201

VERSION := $(shell git describe --tags --always)

.PHONY: dev_server
dev_server:
	go build -o build/dev_server_$(VERSION) ./server
	build/dev_server_$(VERSION) -tcpAddr $(SERVER_TCP_ADDRESS) -udpAddr $(SERVER_UDP_ADDRESS)

.PHONY: dev_app
dev_app:
	go build -ldflags="-X 'main.tcpAddr=$(PUBLIC_TCP_ADDRESS)' -X 'main.udpAddr=$(PUBLIC_UDP_ADDRESS)'" -o build/dev_app_$(VERSION) ./app
	build/dev_app_$(VERSION)

.PHONY: clean
clean:
	-rm -rf ./build

.PHONY: build_server
build_server:
	go build -o build/dev_server_$(VERSION) ./server

.PHONY: build_app
build_app:
	go build -ldflags="-X 'main.tcpAddr=$(PUBLIC_TCP_ADDRESS)' -X 'main.udpAddr=$(PUBLIC_UDP_ADDRESS)'" -o build/dev_app_$(VERSION) ./app
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.tcpAddr=$(PUBLIC_TCP_ADDRESS)' -X 'main.udpAddr=$(PUBLIC_UDP_ADDRESS)'" -o build/dev_app_$(VERSION).exe ./app
