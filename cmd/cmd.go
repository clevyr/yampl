package cmd

import (
	"context"
	"os"

	"gabe565.com/utils/cobrax"
	"github.com/clevyr/yampl/internal/config"
	"github.com/clevyr/yampl/internal/processor"
	"github.com/clevyr/yampl/internal/visitor"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var description = `Yampl (yaml + tmpl) templates YAML values based on line-comments.
YAML data can be piped to stdin or files/dirs can be passed as arguments.

Full reference at ` + termenv.Hyperlink("https://github.com/clevyr/yampl#readme", "github.com/clevyr/yampl")

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "yampl [files | dirs] [-v key=value...]",
		Short:             "Inline YAML templating via line-comments",
		Long:              description,
		DisableAutoGenTag: true,
		ValidArgsFunction: validArgs,
		RunE:              run,
	}
	conf := config.New()
	conf.RegisterFlags(cmd)
	conf.RegisterCompletions(cmd)
	visitor.RegisterCompletion(cmd)
	cmd.SetContext(config.WithContext(context.Background(), conf))

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func validArgs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
}

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if len(args) == 0 {
		if f, ok := cmd.InOrStdin().(*os.File); ok {
			if isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd()) {
				return cmd.Help()
			}
		}
	}

	return processor.Process(conf, args, cmd.InOrStdin(), cmd.OutOrStdout())
}
