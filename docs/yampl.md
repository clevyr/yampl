## yampl

Inline YAML templating via line-comments

```
yampl [-i] [-p prefix] [-v key=value ...] [file ...]
```

### Options

```
  -C, --directory string       dir to hold the generated config (default "./docs")
  -h, --help                   help for yampl
  -i, --inline                 Edit files in-place
      --left-delim string      Override the default left delimiter (default "{{")
  -p, --prefix string          Template prefix. Must begin with '#' (default "#yampl")
      --right-delim string     Override the default right delimiter (default "}}")
  -v, --value stringToString   Define a template variable (default [])
```

