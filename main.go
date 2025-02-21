package main

import (
	"os"

	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/slogx"
	"github.com/clevyr/yampl/cmd"
	"github.com/clevyr/yampl/internal/config"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
