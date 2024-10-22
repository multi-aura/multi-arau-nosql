hello:
	echo "Hello"


run-server:
	cd backend && go run cmd/server/main.go

run-server-websocket:
	cd backend && go run cmd/websocket/main.go

tidy:
	cd backend && go mod tidy

.PHONY: all
all: run-server
