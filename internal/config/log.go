package config

import (
	"io"
	"log/slog"
	"os"
	"time"

	"gabe565.com/utils/slogx"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func (c *Config) InitLog(w io.Writer) {
	InitLog(w, c.LogLevel, c.LogFormat)
}

func InitLog(w io.Writer, level slogx.Level, format slogx.Format) {
	switch format {
	case slogx.FormatJSON:
		slog.SetDefault(slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{
				Level: level,
			}),
		))
	default:
		var color bool
		switch format {
		case slogx.FormatAuto:
			if f, ok := w.(*os.File); ok {
				color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
			}
		case slogx.FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
