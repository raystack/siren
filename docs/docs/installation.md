# Installation

There are several approaches to install Siren CLI

1. [Using a pre-compiled binary](#binary-cross-platform)
2. [Installing with package manager](#macOS)
3. [Installing from source](#building-from-source)
4. [Using the Docker image](#use-the-docker-image)
5. [Using the Helm Chart](#use-the-helm-chart)

#### Binary (Cross-platform)

Download the appropriate version for your platform from [releases](https://github.com/odpf/siren/releases) page. Once downloaded, the binary can be run from anywhere.
You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.
Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

#### macOS

`siren` is available via a Homebrew Tap, and as downloadable binary from the [releases](https://github.com/odpf/siren/releases/latest) page:

```sh
brew install odpf/tap/siren
```

To upgrade to the latest version:

```
brew upgrade siren
```

Check for installed siren version

```sh
siren version
```

#### Linux

`siren` is available as downloadable binaries from the [releases](https://github.com/odpf/siren/releases/latest) page. Download the `.deb` or `.rpm` from the releases page and install with `sudo dpkg -i` and `sudo rpm -i` respectively.

#### Windows

`siren` is available via [scoop](https://scoop.sh/), and as a downloadable binary from the [releases](https://github.com/odpf/siren/releases/latest) page:

```
scoop bucket add siren https://github.com/odpf/scoop-bucket.git
```

To upgrade to the latest version:

```
scoop update siren
```

### Building from source

#### Prerequisites

Siren requires the following dependencies:

- Golang (version 1.18 or above)
- Git

#### Build

Run either of the following commands to clone and compile Siren from source

```sh
git clone git@github.com:odpf/siren.git  (Using SSH Protocol) Or
git clone https://github.com/odpf/siren.git (Using HTTPS Protocol)
```

Install all the golang dependencies

```
make setup
```

Build siren binary file

```
make build
```

Init server config. Customise with your local configurations.

```
make config
```

Run database migrations

```
$ siren server migrate -c config.yaml
```

Start siren server

```
$ siren server start -c config.yaml
```

Initialize client configurations

```
$ siren config init
```

### Use the Docker image

We provide ready to use Docker container images. To pull the latest image:

```
docker pull odpf/siren:latest
```

To pull a specific version:

```
docker pull odpf/siren:v0.4.1
```

### Use the Helm chart

Siren can be installed in Kubernetes using the Helm chart from https://github.com/odpf/charts.

Ensure that the following requirements are met:

- Kubernetes 1.14+
- Helm version 3.x is [installed](https://helm.sh/docs/intro/install/)

Add ODPF chart repository to Helm:

```
helm repo add odpf https://odpf.github.io/charts/
```

You can update the chart repository by running:

```
helm repo update
```

And install it with the helm command line:

```
helm install my-release odpf/siren
```

### Verifying the installation​

To verify if Siren is properly installed, run `siren --help` on your system. You should see help output. If you are executing it from the command line, make sure it is on your PATH or you may get an error about Siren not being found.

```
$ siren --help
```

### Dockerized dependencies

  You will notice there is a [`docker-compose.yaml`](https://github.com/odpf/siren/blob/main/docker-compose.yaml) file contains all dependencies that you need to bootstrap Siren. Inside, it has `postgresql` as a main storage, `cortex ruler` and `cortex alertmanager` as monitoring provider, and `minio` as a backend storage for `cortex`.
