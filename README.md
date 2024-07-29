# Yampl

<img src="./assets/icon.svg" alt="Yampl Icon" width="170" align="right">

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/clevyr/yampl)](https://github.com/clevyr/yampl/releases)
[![Build](https://github.com/clevyr/yampl/actions/workflows/build.yml/badge.svg)](https://github.com/clevyr/yampl/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/clevyr/yampl)](https://goreportcard.com/report/github.com/clevyr/yampl)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clevyr_yampl&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=clevyr_yampl)

Yampl (yaml + tmpl) templates YAML values based on line-comments.

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

  There are two actions available for CI/CD usage:
  - **[clevyr/setup-yampl-action](https://github.com/clevyr/setup-yampl-action):** Installs yampl during a GitHub Action run.
  - **[clevyr/yampl-action](https://github.com/clevyr/yampl-action):** Installs yampl, runs yampl with the given inputs, then optionally creates a commit for you.

</details>

### Docker

<details>
  <summary>Click to expand</summary>

  yampl has a Docker image available at [`ghcr.io/clevyr/yampl`](https://ghcr.io/clevyr/yampl)

  ```shell
  docker pull ghcr.io/clevyr/yampl
  ```

  To use this image, you will need to volume bind the desired directory into the
  Docker container. The container uses `/data` as its workdir, so if you wanted
  to template `example.yaml` in the current directory, you could run:
  ```shell
  docker run --rm -it -v "$PWD:/data" ghcr.io/clevyr/yampl example.yaml ...
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

### AUR (Arch Linux)

<details>
  <summary>Click to expand</summary>

Install [yampl-bin](https://aur.archlinux.org/packages/yampl-bin) with your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.
</details>


## Usage

[View the generated docs](docs/yampl.md) for flag and command reference.
Also, see [templating](#templating) and [example](#examples) sections.

## Examples

### Simple Examples

1. Template with a single value:
    ```shell
    $ cat example.yaml
    name: Clevyr #yampl {{ .name }}
    $ yampl example.yaml -v name='Clevyr Inc.'
    name: Clevyr Inc. #yampl {{ .name }}
    ```

2. Template with multiple values:
    ```shell
    $ cat example.yaml
    image: nginx:stable-alpine #yampl {{ repo current }}:{{ .tag }}
    $ yampl example.yaml -v tag=stable
    image: nginx:stable #yampl {{ repo current }}:{{ .tag }}
    ```

3. Using a [Sprig](https://masterminds.github.io/sprig/) function:
    ```shell
    $ cat example.yaml
    name: Clevyr #yampl {{ upper current }}
    $ yampl example.yaml
    name: CLEVYR #yampl {{ upper current }}
    ```

4. Using the [`repo`](#repo) helper function:
    ```shell
    $ cat example.yaml
    image: nginx:1.20.1 #yampl {{ repo current }}:{{ .tag }}
    $ yampl example.yaml -v tag=1.21.8
    image: nginx:1.21.8 #yampl {{ repo current }}:{{ .tag }}
    ```

### Kubernetes Deployment

<details>
  <summary>Click to expand</summary>

  Here is a simple Kubernetes Deployment with an Nginx image:

  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx
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
            image: nginx:1.26.1 #yampl nginx:{{ .tag }}
            ports:
            - containerPort: 80
  ```

  Notice the yaml comment on the same line as `image`.

  If this file was called `nginx.yaml`, then you could replace the image tag by running:
  ```shell
  yampl -i nginx.yaml -v tag=1.27.0
  ```

  The file would be updated in-place:
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx
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
            image: nginx:1.27.0 #yampl nginx:{{ .tag }}
            ports:
              - containerPort: 80
  ```

  If you wanted to repeat yourself even less, you could use the [`repo`](#repo) function to pull the existing repo through to the output.
  For example, you could change the `image` line to:
  ```yaml
  image: nginx:1.27.0 #yampl {{ repo current }}:{{ .tag }}
  ```

  This would generate the same output, but you wouldn't have to type `nginx` twice.
  This becomes more useful when using custom Docker registries where repo names can get long.

</details>

## Templating

### Variables

All variables passed in with the `-v` flag are available during templating.  
For example, the variable `-v tag=latest` can be used as `{{ .tag }}`.

### Functions

All [Sprig functions](https://masterminds.github.io/sprig/) are available in templates, along with some extras:

#### `current`

Returns the current YAML node's value.

#### `repo`

Splits a Docker repo and tag into the repo component:
```gotemplate
repo "nginx:stable-alpine"
```
The above produces `nginx`.

#### `tag`

Splits a Docker repo and tag into the tag component:
```gotemplate
tag "nginx:stable-alpine"
```
The above produces `stable-alpine`

## Advanced Usage

### Tags

By default, templated values are not explicitly quoted. This can cause
problems with some tools that require specific types. If you require a
specific type for a field, you can add a tag to the template prefix.

Supported tags:

- `#yampl:bool`
- `#yampl:str`
- `#yampl:int`
- `#yampl:float`
- `#yampl:seq`
- `#yampl:map`

For example, the following could be interpreted as either a string or an int:

```shell
$ cat example.yaml
num: #yampl {{ .num }}
$ yampl example.yaml -v num=2009
num: 2009 #yampl {{ .num }}
```

If this field must be a string, you could add the `str` tag:

```shell
$ cat example.yaml
num: #yampl:str {{ .num }}
$ yampl example.yaml -v num=2009
num: "2009" #yampl:str {{ .num }}
```
