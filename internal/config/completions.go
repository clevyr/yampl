package config

import (
	"github.com/clevyr/yampl/internal/util"
	"github.com/spf13/cobra"
)

func (c *Config) RegisterCompletions(cmd *cobra.Command) {
	util.Must(
		cmd.RegisterFlagCompletionFunc(InplaceFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(RecursiveFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(PrefixFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(LeftDelimFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(RightDelimFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(IndentFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FailFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(StripFlag, BoolCompletion),
		cmd.RegisterFlagCompletionFunc(LogLevelFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
			},
		),
		cmd.RegisterFlagCompletionFunc(LogFormatFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{"auto", "color", "plain", "json"}, cobra.ShellCompDirectiveNoFileComp
			},
		),
		cmd.RegisterFlagCompletionFunc(CompletionFlag,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
			},
		),
	)
}

func BoolCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
}
