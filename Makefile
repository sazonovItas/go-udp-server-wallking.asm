EXE = ./cmd/wallking-server/main.go
TEST_EXE = ./tests/client/main.go

CONFIG_PATH = ./config/dev.yaml

setWinEnv:
	set CONFIG_PATH=$(CONFIG_PATH)

setLinuxEnv: ./env.sh
	export CONFIG_PATH=$(CONFIG_PATH)

run: $(EXE)
	go run $(EXE)

test: $(TEST_EXE)
	go run $(TEST_EXE)