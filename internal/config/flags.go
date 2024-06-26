package config

import (
	"github.com/spf13/cobra"
)

const (
	InplaceFlag    = "inplace"
	ValueFlag      = "value"
	ValueFlagShort = "v"
	RecursiveFlag  = "recursive"
	PrefixFlag     = "prefix"
	LeftDelimFlag  = "left-delim"
	RightDelimFlag = "right-delim"
	IndentFlag     = "indent"
	FailFlag       = "fail"
	StripFlag      = "strip"

	LogLevelFlag  = "log-level"
	LogFormatFlag = "log-format"

	CompletionFlag = "completion"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&c.Inplace, InplaceFlag, "i", c.Inplace, "Edit files in place")
	cmd.Flags().StringToStringP(ValueFlag, ValueFlagShort, map[string]string{}, "Define a template variable. Can be used more than once.")
	cmd.Flags().BoolVarP(&c.Recursive, RecursiveFlag, "r", c.Recursive, "Recursively update yaml files in the given directory")
	cmd.Flags().StringVarP(&c.Prefix, PrefixFlag, "p", c.Prefix, "Template comments must begin with this prefix. The beginning '#' is implied.")
	cmd.Flags().StringVar(&c.LeftDelim, LeftDelimFlag, c.LeftDelim, "Override template left delimiter")
	cmd.Flags().StringVar(&c.RightDelim, RightDelimFlag, c.RightDelim, "Override template right delimiter")
	cmd.Flags().IntVarP(&c.Indent, IndentFlag, "I", c.Indent, "Override output indentation")
	cmd.Flags().BoolVarP(&c.Fail, FailFlag, "f", c.Fail, `Exit with an error if a template variable is not set`)
	cmd.Flags().BoolVarP(&c.Strip, StripFlag, "s", c.Strip, "Strip template comments from output")

	cmd.Flags().StringVarP(&c.LogLevel, LogLevelFlag, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	cmd.Flags().StringVar(&c.LogFormat, LogFormatFlag, c.LogFormat, "Log format (auto, color, plain, json)")

	cmd.Flags().StringVar(&c.Completion, CompletionFlag, c.Completion, "Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.")
}
