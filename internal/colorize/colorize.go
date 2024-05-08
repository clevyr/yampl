package colorize

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
	"github.com/mattn/go-isatty"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
	return escape + "[" + strconv.Itoa(int(attr)) + "m"
}

func shouldColor(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	if f, ok := w.(*os.File); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}
	return false
}

func WriteString(w io.Writer, s string) error {
	if shouldColor(w) {
		s = Colorize(s)
	}

	_, err := io.WriteString(w, s)
	return err
}

func Colorize(s string) string {
	// https://github.com/mikefarah/yq/blob/v4.43.1/pkg/yqlib/color_print.go
	tokens := lexer.Tokenize(s)
	var p printer.Printer
	p.Bool = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.Number = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiMagenta),
			Suffix: format(color.Reset),
		}
	}
	p.MapKey = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgCyan),
			Suffix: format(color.Reset),
		}
	}
	p.Anchor = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.Alias = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgHiYellow),
			Suffix: format(color.Reset),
		}
	}
	p.String = func() *printer.Property {
		return &printer.Property{
			Prefix: format(color.FgGreen),
			Suffix: format(color.Reset),
		}
	}
	s = p.PrintTokens(tokens)
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}
