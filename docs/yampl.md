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

Simple Examples (Full example at https://github.com/clevyr/go-yampl#examples):

 $ echo 'name: Clevyr #yampl {{ .name }}' | yampl -v name='Clevyr Inc.'  
 name: Clevyr Inc. #yampl {{ .name }}  
 $ echo 'name: Clevyr #yampl {{ upper .Value }}' | yampl  
 name: CLEVYR #yampl {{ upper .Value }}  
 $ echo 'image: nginx:stable-alpine #yampl {{ repo .Value }}:{{ .tag }}' | yampl -v tag=stable  
 image: nginx:stable #yampl {{ repo .Value }}:{{ .tag }}

Template Function Reference:
 - https://github.com/clevyr/go-yampl#functions
 - https://masterminds.github.io/sprig/

Template Variable Reference:
 - https://github.com/clevyr/go-yampl#variables


```
yampl [-i] [-p prefix] [-v key=value ...] [file ...]
```

### Options

```
      --completion string      Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.
  -h, --help                   help for yampl
  -I, --indent int             Override output indentation (default 2)
  -i, --inplace                Update files inplace
      --left-delim string      Override the left delimiter (default "{{")
  -p, --prefix string          Line-comments are ignored unless this prefix is found. Prefix must begin with '#' (default "#yampl")
      --right-delim string     Override the right delimiter (default "}}")
  -s, --strict                 Trigger an error if a template variable is missing
  -v, --value stringToString   Define a template variable. Can be used more than once. (default [])
```

