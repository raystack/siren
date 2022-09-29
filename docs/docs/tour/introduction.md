# Introduction

This tour introduces you to how to use Siren to manage alerting rules and send notifications. We will integrate Siren with a monitoring provider, configuring alerting rules for the provider, accepting triggered alerts from the provider, sending notification from the triggered alerts, and sending notification on-demand. The tour takes approximately 20 minutes to complete.

This tour will mainly use `siren` CLI. To see what commands does `siren` CLI has, you can see the help for a command using `--help`:

```shell
$ siren --help
$ siren provider --help
```

## Prerequisites
- Install `git`
- Install `go 1.16` 
- Install `make`
- Install `docker`
- Install `docker-compose`

## Clone the Git repository

First, clone the Git repository of siren. From the development directory of your choice, run this command:

```shell
$ git clone https://github.com/odpf/siren.git
```

You will notice there is a `docker-compose.yaml` file contains all dependencies that you need to bootstrap Siren. Inside, it has `postgresql` as a main storage, `cortex ruler` and `cortex alertmanager` as monitoring provider, and `minio` as a backend storage for `cortex`.