package main

import (
	"asm-game/server/cmd/wallking-server/server"
	"log"
	"log/slog"
	"os"

	config "asm-game/server/internal/config"
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

	log := setupLogger(cfg.Env)

	sv := server.New(log, cfg.Address, cfg.Port)
	sv.Up()
	defer sv.Down()

	buf := make([]byte, 256)
	for {
		n, addr, err := sv.ListenCon.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error to read from the socker", slog.String("error", err.Error()))
		}

		// log.Debug(
		// 	"Message",
		// 	slog.String("address", addr.String()),
		// 	slog.String("msg", string(buf[:n])),
		// )
		go sv.UpdatePlConn(addr.String(), buf[:n])
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
