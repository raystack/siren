# Start Server

## 1. Build Siren

We'll start our tour by building siren and generating the default configuration. Run this command to build `siren`.

```shell
$ make build
```

Once succeed, you will get this message and there is a `siren` binary in the folder.

```
 > building siren version v0.4.1-next
go build -ldflags "-X main.Version="v0.4.1-next"" "github.com/odpf/siren"
 - build complete
```

## 2. Generate Siren config

Before starting Siren server, you need to have a proper config file. Run server config initialization command to generate a config file with the defaults value.

```shell
$ ./siren server init
```

You will have a new config file `./config.yaml` with the default values. The default values would be the minimum config to get Siren up and running.

## 3. Start-up Siren dependencies

We are using `docker-compose` to start-up all dependencies of Siren. See [introduction](./introduction.md) to see what components are defined in `docker-compose`.

To start all dependencies of Siren, run.

```shell
docker-compose up -d
```

If you haven't got dependencies (postgresql, cortex, etc) images in your local, this might take a while to fetch all images. If all dependencies are successfully started, you will get this message in the terminal.

```shell
[+] Running 9/9
 ⠿ Network siren_siren              Created 0.0s
 ⠿ Network siren_default            Created 0.0s
 ⠿ Volume "siren_data1-1"           Created 0.0s
 ⠿ Volume "siren_siren_dbdata"      Created 0.0s
 ⠿ Container siren_postgres         Started 0.7s
 ⠿ Container siren-minio1-1         Started 0.7s
 ⠿ Container siren_cortex_am        Started 1.3s
 ⠿ Container siren-createbuckets-1  Started 1.1s
 ⠿ Container siren_cortex_all       Started 1.6s
```

## 4. Run Siren server

If all dependencies are up and running, you can start running Siren server. To run a Siren server, we could just need to run a command.

```shell
$ ./siren server start
```

Siren will auto-recognize `./config.yaml` file as its config file and fetch all configs inside it to be used in Siren server. The default config also runs two notification workers inside Siren server. If it is running properly, you will see this logs inside terminal.

```shell
{"severity":"INFO","timestamp":"2022-10-20T11:19:46.486792+07:00","caller":"log/zap.go:24","message":"running worker","serviceContext":{"service":"siren"},"id":"a2ef71b6-b8a6-4d3c-8358-bfe406f268d6"}
{"severity":"INFO","timestamp":"2022-10-20T11:19:46.486766+07:00","caller":"log/zap.go:24","message":"running worker","serviceContext":{"service":"siren"},"id":"361945b4-aaf3-4e45-afb9-605888c1ebda"}
{"severity":"INFO","timestamp":"2022-10-20T11:19:46.487412+07:00","caller":"log/zap.go:24","message":"server is running","serviceContext":{"service":"siren"},"host":"localhost","port":8080}
```

## 5. Migrate Siren DB

To migrate Siren's DB (Postgres), you just need to run this command.

```shell
$ ./siren server migrate
```

This will create all necessary tables inside Siren's DB.

## 6. Siren client CLI Config Initialization

Siren client CLI require some configurations to communicate to Siren server. One of them is Siren host. To create Siren client config, run.

```shell
./siren config init
```

This will autogenerate Siren client config inside path `${HOME}/odpf/siren.yaml` with the default value. Ideally, you don't need to do anything with the value, the default value should be compatible in this tour.

Once Siren server is up and running, we can start configuring Siren to our needs.
