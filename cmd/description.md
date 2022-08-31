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
