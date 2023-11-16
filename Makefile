SHELL=/usr/bin/env bash

run: ./env.sh ./cmd/wallking-server/main.go
	source ./env.sh && go run ./cmd/wallking-server/main.go
