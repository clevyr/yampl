## yampl

Inline YAML templating via line-comments

### Synopsis

Yampl (yaml + tmpl) templates YAML values based on line-comments.
YAML data can be piped to stdin or files/dirs can be passed as arguments.

```
yampl [files | dirs] [-v key=value...] [flags]
```

### Options

```
  -h, --help                     help for yampl
      --ignore-template-errors   Continue processing a file even if a template fails
      --ignore-unset-errors      Exit with an error if a template variable is not set (default true)
  -I, --indent int               Override output indentation (default 2)
  -i, --inplace                  Edit files in place
      --left-delim string        Override template left delimiter (default "{{")
      --log-format string        Log format (one of auto, color, plain, json) (default "auto")
  -l, --log-level string         Log level (one of trace, debug, info, warn, error) (default "info")
      --no-source-comment        Disables source path comment when run against multiple files or a dir
  -p, --prefix string            Template comments must begin with this prefix. The beginning '#' is implied. (default "#yampl")
      --right-delim string       Override template right delimiter (default "}}")
  -s, --strip                    Strip template comments from output
  -v, --var stringToString       Define a template variable. Can be used more than once. (default [])
      --version                  version for yampl
```

