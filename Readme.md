# Lacks Chat Server

## Overview

A chat server writen in Golang that supports both TCP and HTTP (partial) operations.

## Prerequisites

- Golang (https://golang.org/)
- dep (https://golang.github.io/dep/docs/installation.html)

## Configuration

| key       | default value |
| --------- | ------------- |
| TCPPort   | 3000          |
| HTTPPort  | 8000          |
| HostAddr  | localhost     |
| LogOutput | lacks.log     |

Utilizes the [viper](https://github.com/spf13/viper) package for configuration, so you can use `yaml`, `toml`, or `json` for the config

### Example YAML Config (sane defaults)

```yaml
TCPPort:   "3000"
HTTPPort:  "4000"
HostAddr:  "localhost"
LogOutput: "lacks.log"
```

Configuration data is read from the following locations, and in the following order (Last In Wins).

- `/etc/chat-lacks/config.{yaml,json,toml}`
- `$HOME/.chat-lacks/config.{yaml,json.toml}`
- `./config.{yaml,json.toml}`

## Dependencies

- [viper](https://github.com/spf13/viper)

## Install Missing 3rd Party Dependencies (Vendor)

Clone the repo and ensure deps are installed:

```bash
dep ensure
```

## Run the server

do 1 of the following...

```bash
go run main.go
```

or

```bash
make && ./chat-lacks
```

or

```bash
go build && ./chat-lacks
```

or

```bash
go build && go install && chat-lacks
```

### Chatroom Useage

Users can interact with the chatroom through TCP and/or HTTP connections.

#### TCP Connections

```bash
telnet localhost 3000
```

#### HTTP Connections

Http Connections use a cookie (Not secure) to maintain a user session.

```bash
curl localhost:4000/
```

or

```bash
curl --cookie "username=troy" -H "Content-Type: application/json" localhost:4000
```

**NOTE: The HTTP API is incomplete**

## Roadmap

- Finish up the HTTP Api and add in a WebSocket endpoint
- Create a web app that utilizes the /ws endpoint
- More testing
