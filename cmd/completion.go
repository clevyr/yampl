package cmd

import (
	"errors"
	"fmt"

	"github.com/clevyr/yampl/internal/config"
	"github.com/spf13/cobra"
)

var ErrInvalidShell = errors.New("invalid shell")

func completion(cmd *cobra.Command, shell string) error {
	switch shell {
	case config.Bash:
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	case config.Zsh:
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	case config.Fish:
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	case config.Powershell:
		return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
	default:
		return fmt.Errorf("%w: %s", ErrInvalidShell, shell)
	}
}
