package config

import (
	"fmt"
	"io"

	"github.com/clevyr/yampl/internal/colorize"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func logLevel(level string) zerolog.Level {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil || parsedLevel == zerolog.NoLevel {
		if level == "warning" {
			parsedLevel = zerolog.WarnLevel
		} else {
			log.Warn().Str("value", level).Msg("invalid log level. defaulting to info.")
			parsedLevel = zerolog.InfoLevel
		}
	}
	return parsedLevel
}

func logFormat(out io.Writer, format string) io.Writer {
	switch format {
	case JSON:
		return out
	default:
		sprintf := fmt.Sprintf
		var useColor bool
		switch format {
		case Auto:
			if useColor = colorize.ShouldColor(out); !useColor {
				break
			}
			fallthrough
		case Color:
			useColor = true
			color.NoColor = false
			sprintf = color.New(color.Bold).Sprintf
		case Plain:
		default:
			log.Warn().Str("value", format).Msg("invalid log formatter. defaulting to auto.")
			return logFormat(out, Auto)
		}

		return zerolog.ConsoleWriter{
			Out:     out,
			NoColor: !useColor,
			FormatMessage: func(i interface{}) string {
				return sprintf("%-45s", i)
			},
		}
	}
}

func initLog(cmd *cobra.Command) {
	level, err := cmd.Flags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(logLevel(level))

	format, err := cmd.Flags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(logFormat(cmd.ErrOrStderr(), format))
}
