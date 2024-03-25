package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func registerLogFlags(cmd *cobra.Command) {
	var err error

	cmd.Flags().StringP("log-level", "l", "info", "Log level (trace, debug, info, warning, error, fatal, panic)")
	err = cmd.RegisterFlagCompletionFunc(
		"log-level",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		panic(err)
	}

	cmd.Flags().String("log-format", "color", "Log format (auto, color, plain, json)")
	err = cmd.RegisterFlagCompletionFunc(
		"log-format",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"auto", "color", "plain", "json"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		panic(err)
	}
}

func initLogLevel(level string) log.Level {
	parsed, err := log.ParseLevel(level)
	if err != nil {
		log.WithField("level", level).Warn("invalid log level. defaulting to info.")
		parsed = log.InfoLevel
	}
	log.SetLevel(parsed)

	return parsed
}

//nolint:ireturn
func initLogFormat(format string) log.Formatter {
	var formatter log.Formatter = &log.TextFormatter{}
	switch format {
	case "auto", "a":
		break
	case "color", "c":
		formatter.(*log.TextFormatter).ForceColors = true
	case "plain", "p":
		formatter.(*log.TextFormatter).DisableColors = true
	case "json", "j":
		formatter = &log.JSONFormatter{}
	default:
		log.WithField("format", format).Warn("invalid log formatter. defaulting to auto.")
	}
	log.SetFormatter(formatter)
	return formatter
}

func initLog(cmd *cobra.Command) {
	level, err := cmd.Flags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	initLogLevel(level)

	format, err := cmd.Flags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	initLogFormat(format)
}
