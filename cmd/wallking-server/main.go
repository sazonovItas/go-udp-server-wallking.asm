package main

import (
	"asm-game/server/cmd/wallking-server/server"
	config "asm-game/server/internal/config"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error occured: %s", err.Error())
	}
	_ = cfg

	logger := setupLogger(cfg.Env)

	sv := server.New(logger, cfg.Address, cfg.Port)
	sv.Up()
	defer sv.Down()

	go sv.SendMessages()

	for {
		buf := make([]byte, 256)
		n, addr, err := sv.ListenCon.ReadFromUDP(buf)
		if err != nil {
			logger.Error("Error to read from the socket", slog.String("error", err.Error()))
		}

		//	logger.Debug("message from socket", slog.String("msg", string(buf[:n])))

		go sv.HandleMsg(addr.String(), buf[:n])
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
