package main

import (
	"log/slog"
	"os"

	"gabe565.com/utils/cobrax"
	"github.com/clevyr/yampl/cmd"
	"github.com/clevyr/yampl/internal/config"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
