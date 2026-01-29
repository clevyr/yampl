package config

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"

	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

const (
	VarFlag = "var"

	InplaceFlag         = "inplace"
	PrefixFlag          = "prefix"
	LeftDelimFlag       = "left-delim"
	RightDelimFlag      = "right-delim"
	IndentFlag          = "indent"
	StripFlag           = "strip"
	NoSourceCommentFlag = "no-source-comment"

	IgnoreUnsetErrorsFlag    = "ignore-unset-errors"
	IgnoreTemplateErrorsFlag = "ignore-template-errors"

	LogLevelFlag  = "log-level"
	LogFormatFlag = "log-format"

	// Deprecated: Replaced by VarFlag.
	ValueFlag = "value"
	// Deprecated: Removed. Yampl will always recurse if a given path is a directory.
	RecursiveFlag = "recursive"
	// Deprecated: Replaced by IgnoreUnsetErrorsFlag and IgnoreTemplateErrorsFlag.
	FailFlag = "fail"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().Var(c.valuesStringToString, ValueFlag, "Define a template variable. Can be used more than once.")

	cmd.Flags().BoolVarP(&c.Inplace, InplaceFlag, "i", c.Inplace, "Edit files in place")
	cmd.Flags().StringVarP(&c.Prefix, PrefixFlag, "p", c.Prefix,
		"Template comments must begin with this prefix. The beginning '#' is implied.",
	)
	cmd.Flags().StringVar(&c.LeftDelim, LeftDelimFlag, c.LeftDelim, "Override template left delimiter")
	cmd.Flags().StringVar(&c.RightDelim, RightDelimFlag, c.RightDelim, "Override template right delimiter")
	cmd.Flags().IntVarP(&c.Indent, IndentFlag, "I", c.Indent, "Override output indentation")
	cmd.Flags().BoolVarP(&c.Strip, StripFlag, "s", c.Strip, "Strip template comments from output")
	cmd.Flags().BoolVar(&c.NoSourceComment, NoSourceCommentFlag, c.NoSourceComment,
		"Disables source path comment when run against multiple files or a dir",
	)

	cmd.Flags().BoolVar(&c.IgnoreUnsetErrors, IgnoreUnsetErrorsFlag, c.IgnoreUnsetErrors,
		"Exit with an error if a template variable is not set",
	)
	cmd.Flags().BoolVar(&c.IgnoreTemplateErrors, IgnoreTemplateErrorsFlag, c.IgnoreTemplateErrors,
		"Continue processing a file even if a template fails",
	)

	cmd.Flags().VarP(&c.LogLevel, LogLevelFlag, "l", "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	cmd.Flags().Var(&c.LogFormat, LogFormatFlag, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	// Deprecated
	cmd.Flags().VarP(c.valuesStringToString, VarFlag, "v", "Define a template variable. Can be used more than once.")
	cmd.Flags().BoolP(RecursiveFlag, "r", true, "Recursively update yaml files in the given directory")
	cmd.Flags().BoolP(FailFlag, "f", false, `Exit with an error if a template variable is not set`)
	if err := errors.Join(
		cmd.Flags().MarkDeprecated(ValueFlag, "use --"+VarFlag+" instead"),
		cmd.Flags().MarkDeprecated(RecursiveFlag, cmd.Name()+" will always recurse if a given path is a directory"),
		cmd.Flags().MarkDeprecated(FailFlag,
			"use --"+IgnoreUnsetErrorsFlag+" and --"+IgnoreTemplateErrorsFlag+" instead",
		),
	); err != nil {
		panic(err)
	}

	c.InitLog(cmd.ErrOrStderr())
	cmd.Flags().SetOutput(DeprecatedWriter{})
}

type DeprecatedWriter struct{}

func (d DeprecatedWriter) Write(b []byte) (int, error) {
	slog.Warn(string(bytes.TrimSpace(b)))
	return len(b), nil
}
