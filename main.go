package main

import (
	"github.com/clevyr/go-yampl/cmd"
	"os"
)

//go:generate go run internal/cmd/docs/main.go --directory=docs

func main() {
	if err := cmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
