package main

import (
	"os"

	"github.com/clevyr/yampl/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
