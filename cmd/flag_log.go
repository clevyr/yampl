package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel  string
	logFormat string
)

func init() {
	var err error

	Command.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (trace, debug, info, warning, error, fatal, panic)")
	err = Command.RegisterFlagCompletionFunc(
		"log-level",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		panic(err)
	}

	Command.Flags().StringVar(&logFormat, "log-format", "color", "Log format (auto, color, plain, json)")
	err = Command.RegisterFlagCompletionFunc(
		"log-format",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"auto", "color", "plain", "json"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initLog)
}

func initLogLevel(level string) log.Level {
	parsed, err := log.ParseLevel(level)
	if err != nil {
		log.WithField("level", logLevel).Warn("invalid log level. defaulting to info.")
		logLevel = "info"
		parsed = log.InfoLevel
	}
	log.SetLevel(parsed)

	return parsed
}

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
		log.WithField("format", logFormat).Warn("invalid log formatter. defaulting to auto.")
	}
	log.SetFormatter(formatter)
	return formatter
}

func initLog() {
	initLogLevel(logLevel)
	initLogFormat(logFormat)
}
