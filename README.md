# YAMPL

YAMPL (yaml + tmpl) is a simple templating tool that can template yaml values based on line-comments.

## Installation

YAMPL is available in brew or as a Docker container.

#### Brew

```shell
brew install clevyr/tap/yampl
```

#### Docker

```shell
docker pull ghcr.io/clevyr/yampl
```

## Usage

[View the help docs for usage information.](docs/yampl.md)

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
repo "nginx:stable-alpine"
```
The above produces `stable-alpine`

### Variables

All variables passed in with the `-v` flag are available during templating.
For example, a variable given as `-v tag=latest` can be used as `{{ .tag }}`.

The previous value is always available via `.Value` (`.Val` or `.V` if you're feeling lazy).

### Example

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
