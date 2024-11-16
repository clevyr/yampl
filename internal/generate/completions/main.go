package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"gabe565.com/utils/cobrax"
	"github.com/clevyr/yampl/cmd"
)

func main() {
	if err := os.RemoveAll("completions"); err != nil {
		panic(err)
	}

	rootCmd := cmd.New()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	for _, shell := range []cobrax.Shell{cobrax.Bash, cobrax.Zsh, cobrax.Fish} {
		if err := cobrax.GenCompletion(rootCmd, shell); err != nil {
			panic(err)
		}

		path := filepath.Join("completions", string(shell))
		if err := os.MkdirAll(path, 0o777); err != nil {
			panic(err)
		}

		switch shell {
		case cobrax.Bash:
			path = filepath.Join(path, rootCmd.Name())
		case cobrax.Zsh:
			path = filepath.Join(path, "_"+rootCmd.Name())
		case cobrax.Fish:
			path = filepath.Join(path, rootCmd.Name()+".fish")
		}

		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(f, &buf); err != nil {
			panic(err)
		}

		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}
