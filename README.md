# Yampl

[![Go](https://github.com/clevyr/go-yampl/actions/workflows/go.yml/badge.svg)](https://github.com/clevyr/go-yampl/actions/workflows/go.yml)

Yampl (yaml + tmpl) is a simple tool to template yaml values based on line-comments.

## Installation

Yampl is available in brew or as a Docker container.

### Homebrew (macOS, Linux)

<details>
  <summary>Click to expand</summary>

  ```shell
  brew install clevyr/tap/yampl
  ```
</details>

### GitHub Actions

<details>
  <summary>Click to expand</summary>

  There is an Action provided for use during CI/CD. See [clevyr/yampl-action](https://github.com/clevyr/yampl-action) for more details.
</details>

### Docker

<details>
  <summary>Click to expand</summary>

  yampl has a Docker image available at `ghcr.io/clevyr/yampl`

  ```shell
  docker pull ghcr.io/clevyr/yampl
  ```

  To use this image, you will need to volume bind the desired directory into the
  Docker container. The container uses `/data` as its workdir, so if you wanted
  to template `example.yaml` in the current directory, you could run:
  ```shell
  docker run --rm -it -v "$PWD:/data" ghcr.io/clevyr/yampl yampl example.yaml ...
  ```
</details>

### APT Repository (Ubuntu, Debian)

<details>
  <summary>Click to expand</summary>

  1. If you don't have it already, install the `ca-certificates` package
     ```shell
     sudo apt install ca-certificates
     ```

  2. Add Clevyr's apt repository
     ```
     echo 'deb [trusted=yes] https://apt.clevyr.com /' | sudo tee /etc/apt/sources.list.d/clevyr.list
     ```

  3. Update apt repositories
     ```shell
     sudo apt update
     ```

  4. Install yampl
     ```shell
     sudo apt install yampl
     ```
</details>

### RPM Repository (CentOS, RHEL)

<details>
  <summary>Click to expand</summary>

  1. If you don't have it already, install the `ca-certificates` package
     ```shell
     sudo yum install ca-certificates
     ```

  2. Add Clevyr's rpm repository to `/etc/yum.repos.d/clevyr.repo`
     ```ini
     [clevyr]
     name=Clevyr
     baseurl=https://rpm.clevyr.com
     enabled=1
     gpgcheck=0
     ```

  3. Install yampl
     ```shell
     sudo yum install yampl
     ```
</details>

## Templating

### Functions

All [Sprig functions](https://masterminds.github.io/sprig/) are available in templates, along with some extra functions:

### `repo`

Splits a Docker repo and tag into the repo component:
```gotemplate
repo "nginx:stable-alpine"
```
The above produces `nginx`.

### `tag`

Splits a Docker repo and tag into the tag component:
```gotemplate
tag "nginx:stable-alpine"
```
The above produces `stable-alpine`

### Variables

All variables passed in with the `-v` flag are available during templating.  
For example, the variable `-v tag=latest` can be used as `{{ .tag }}`.

The previous value is always available via `.Value` (`.Val` or `.V` if you're feeling lazy).

## Examples

### Simple Examples

```shell
$ echo 'name: Clevyr #yampl {{ .name }}' | yampl -v name='Clevyr Inc.'
name: Clevyr Inc. #yampl {{ .name }}
$ echo 'name: Clevyr #yampl {{ upper .Value }}' | yampl
name: CLEVYR #yampl {{ upper .Value }}
$ echo 'image: nginx:stable-alpine #yampl {{ repo .Value }}:{{ .tag }}' | yampl -v tag=stable
image: nginx:stable #yampl {{ repo .Value }}:{{ .tag }}
```

### Full Example

Here is a simple Kubernetes nginx Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.20.2 #yampl nginx:{{ .tag }}
          ports:
          - containerPort: 80
```

In this example, notice the yaml comment to the right of the `image`.

If this file was called `nginx.yaml`, then we could replace the image tag by running the following:
```shell
yampl -i nginx.yaml -v tag=1.21.6
```

The file would be updated in-place and would end up looking like:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.21.6 #yampl nginx:{{ .tag }}
          ports:
            - containerPort: 80
```

If I wanted to repeat myself even less, I could utilize the `repo` function to pull the existing repo through.
I could define the `image` template as:
```yaml
image: nginx:1.21.6 #yampl {{ repo .Value }}:{{ .tag }}
```

This would generate the same output, but I didn't have to type `nginx` twice.
This becomes more useful when using custom Docker registries where repo names can get quite long.

## Tags

By default, templated values are not explicitly quoted. This can cause
problems with some tools that require specific types. If you require a
field specific type for a field, you can add a tag to the template prefix.

Supported tags:

- `#yampl:bool`
- `#yampl:str`
- `#yampl:int`
- `#yampl:float`
- `#yampl:seq`
- `#yampl:map`

For example, the following could be interpreted as either a string or an int:

```shell
$ echo 'num: "" #yampl {{ .num }}' | yampl -v num=2009
num: 2009 #yampl {{ .num }}
```

If this field must be a string, you could add the `str` tag:

```shell
$ echo 'num: "" #yampl:str {{ .num }}' | yampl -v num=2009
num: "2009" #yampl:str {{ .num }}
```


## Usage

[View the generated docs for usage information.](docs/yampl.md)
