# problematic-api-server

> I got 429 problems but an API ain't one.

![cover](/assets/cover.jpg "cover")

A server that has examples for every problematic API consumption issue one may
face when consuming APIs.

This project is mostly for educational purposes, but it can be used as a
reference for how to handle different API issues.

## Usage

```sh
go run cmd/server.go
```

### Configuration

The server can be configured using environment variables, with the `PROBLEMATIC`
prefix:

* `PROBLEMATIC_HOST`: The host to listen on. Defaults to `0.0.0.0`.
* `PROBLEMATIC_PORT`: The port to listen on. Defaults to `4578`.
* `PROBLEMATIC_LOG_LEVEL`: The log level to use. Defaults to `debug`. Use one
  of the following: panic, fatal, error, warn, info, debug, trace.
* `PROBLEMATIC_LOG_FORMAT`: The log format to use. Defaults to `text`.
  * Use `text` for human-readable logs.
  * Use `json` for machine-readable logs. Useful when running in a deployed context.

For example, to change the port, run:

```sh
PROBLEMATIC_PORT=2345 go run cmd/server.go
```

> To see all available configuration options, please review `cmd/config.go`.

## Development

### Pre-commit

This project uses [pre-commit](https://pre-commit.com/) to run some checks
before committing. To install it, run:

```sh
brew install pre-commit
pre-commit install
```

### API

This project is based mostly around the API spec in `spec/openapi.yaml`. Use
`./scripts/generate.sh` to generate the server and client if you change the API.

> Note: some files are ignored when generating the server and client. Make sure
> that your changes propagate to the ignored files as well.

To implement a specific service, go to `/server/go/XXX_service.go` and
implement the methods there.
