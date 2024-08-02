package config

import (
	"bytes"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	VarFlag = "var"

	InplaceFlag    = "inplace"
	PrefixFlag     = "prefix"
	LeftDelimFlag  = "left-delim"
	RightDelimFlag = "right-delim"
	IndentFlag     = "indent"
	StripFlag      = "strip"

	IgnoreUnsetErrorsFlag    = "ignore-unset-errors"
	IgnoreTemplateErrorsFlag = "ignore-template-errors"

	LogLevelFlag  = "log-level"
	LogFormatFlag = "log-format"

	CompletionFlag = "completion"

	ValueFlag = "value"
	RecursiveFlag = "recursive"
	FailFlag = "fail"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().Var(c.valuesStringToString, ValueFlag, "Define a template variable. Can be used more than once.")

	cmd.Flags().BoolVarP(&c.Inplace, InplaceFlag, "i", c.Inplace, "Edit files in place")
	cmd.Flags().StringVarP(&c.Prefix, PrefixFlag, "p", c.Prefix, "Template comments must begin with this prefix. The beginning '#' is implied.")
	cmd.Flags().StringVar(&c.LeftDelim, LeftDelimFlag, c.LeftDelim, "Override template left delimiter")
	cmd.Flags().StringVar(&c.RightDelim, RightDelimFlag, c.RightDelim, "Override template right delimiter")
	cmd.Flags().IntVarP(&c.Indent, IndentFlag, "I", c.Indent, "Override output indentation")
	cmd.Flags().BoolVarP(&c.Strip, StripFlag, "s", c.Strip, "Strip template comments from output")

	cmd.Flags().BoolVar(&c.IgnoreUnsetErrors, IgnoreUnsetErrorsFlag, c.IgnoreUnsetErrors, "Exit with an error if a template variable is not set")
	cmd.Flags().BoolVar(&c.IgnoreTemplateErrors, IgnoreTemplateErrorsFlag, c.IgnoreTemplateErrors, "Continue processing a file even if a template fails")

	cmd.Flags().StringVarP(&c.LogLevel, LogLevelFlag, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	cmd.Flags().StringVar(&c.LogFormat, LogFormatFlag, c.LogFormat, "Log format (auto, color, plain, json)")

	cmd.Flags().StringVar(&c.Completion, CompletionFlag, c.Completion, "Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.")

	// Deprecated
	cmd.Flags().VarP(c.valuesStringToString, VarFlag, "v", "Define a template variable. Can be used more than once.")
	cmd.Flags().BoolP(RecursiveFlag, "r", true, "Recursively update yaml files in the given directory")
	cmd.Flags().BoolP(FailFlag, "f", false, `Exit with an error if a template variable is not set`)
	if err := errors.Join(
		cmd.Flags().MarkDeprecated(ValueFlag, "use --"+VarFlag+" instead"),
		cmd.Flags().MarkDeprecated(RecursiveFlag, cmd.Name()+" will always recurse if a given path is a directory"),
		cmd.Flags().MarkDeprecated(FailFlag, "use --"+IgnoreUnsetErrorsFlag+" and --"+IgnoreTemplateErrorsFlag+" instead"),
	); err != nil {
		panic(err)
	}

	initLog(cmd)
	cmd.Flags().SetOutput(DeprecatedWriter{})
}

type DeprecatedWriter struct{}

func (d DeprecatedWriter) Write(b []byte) (int, error) {
	log.Warn().Msg(string(bytes.TrimSpace(b)))
	return len(b), nil
}
