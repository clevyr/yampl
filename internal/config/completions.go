package config

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	Bash       = "bash"
	Zsh        = "zsh"
	Fish       = "fish"
	Powershell = "powershell"
)

func (c *Config) RegisterCompletions(cmd *cobra.Command) {
	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(InplaceFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(PrefixFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(LeftDelimFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(RightDelimFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(IndentFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(IgnoreUnsetErrorsFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(IgnoreTemplateErrorsFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(StripFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(LogLevelFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return LogLevelStrings(), cobra.ShellCompDirectiveNoFileComp
			},
		),
		cmd.RegisterFlagCompletionFunc(LogFormatFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return LogFormatStrings(), cobra.ShellCompDirectiveNoFileComp
			},
		),
		cmd.RegisterFlagCompletionFunc(CompletionFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{Bash, Zsh, Fish, Powershell}, cobra.ShellCompDirectiveNoFileComp
			},
		),
	); err != nil {
		panic(err)
	}
}

func BoolCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
}
