## yampl

Inline YAML templating via line-comments

### Synopsis

Yampl (yaml + tmpl) is a simple tool to template yaml values based on line-comments.

This command can work on stdin/stdout or on files similarly to GNU sed:
 - If no positional args are given, it will listen for input on stdin.
 - It will continue to read input until EOF is encountered.
 - By default, it will print to stdout.
 - The `-i` flag will make it update files in-place.
 - Multiple files can be given, and they will all be templated.

Simple Examples (Full example at https://github.com/clevyr/yampl#examples):

 $ echo 'name: Clevyr #yampl {{ .name }}' | yampl -v name='Clevyr Inc.'  
 name: Clevyr Inc. #yampl {{ .name }}  
 $ echo 'name: Clevyr #yampl {{ upper .Value }}' | yampl  
 name: CLEVYR #yampl {{ upper .Value }}  
 $ echo 'image: nginx:stable-alpine #yampl {{ repo .Value }}:{{ .tag }}' | yampl -v tag=stable  
 image: nginx:stable #yampl {{ repo .Value }}:{{ .tag }}

Template Function Reference:
 - https://github.com/clevyr/yampl#functions
 - https://masterminds.github.io/sprig/

Template Variable Reference:
 - https://github.com/clevyr/yampl#variables


```
yampl [-i] [-p prefix] [-v key=value ...] [file ...]
```

### Options

```
      --completion string      Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.
  -f, --fail                   Exit with an error if a template variable is not set
  -h, --help                   help for yampl
  -I, --indent int             Override output indentation (default 2)
  -i, --inplace                Edit files in place
      --left-delim string      Override template left delimiter (default "{{")
      --log-format string      Log format (auto, color, plain, json) (default "color")
  -l, --log-level string       Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
  -p, --prefix string          Template comments must begin with this prefix. The beginning '#' is implied. (default "#yampl")
  -r, --recursive              Recursively update yaml files in the given directory
      --right-delim string     Override template right delimiter (default "}}")
  -s, --strip                  Strip template comments from output
  -v, --value stringToString   Define a template variable. Can be used more than once. (default [])
```

