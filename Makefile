EXE=./cmd/wallking-server/main.go
TEST_EXE=./tests/client/main.go

CONFIG_PATH=./config/dev.yaml

run: $(EXE)
	go run $(EXE)

build: $(EXE)
	go build -o WallKingServer.exe $(EXE)

test: $(TEST_EXE)
	go run $(TEST_EXE)