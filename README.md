# byrate

A fully self-contained, portable, self-hosted internet download speed test tool written in Go. Easily measure download performance without any external dependencies or third-party services—just run the binary and test.

![Banner](./banner.png)

## Download

```bash
curl -sL https://bit.ly/byrate-dl | sh
```

## Get started

Download the appropriate binary file for your OS from [releases](https://github.com/dipakw/byrate/releases) or the commands above.

| Description                               | Example                                |
|-------------------------------------------|----------------------------------------|
| Default start (starts on `[::1]:14000`)   | `byrate`                               |
| Custom host / port                        | `byrate s -h=localhost -p=15000`       |
| On a Unix socket                          | `byrate s -u -h=/tmp/byrate.sock`      |

## CLI Usage

```
Usage:
  byrate <command> [options]

Commands:
  version, v   Show version
  start, s     Start the server (default)
  help, h      Show this help message

Options:
  --host, -h    Server host (default: ::1)
  --port, -p    Server port (default: 14000)
  --unix, -u    Use unix socket instead of TCP

Notes:
  - All options can use either --long or -short forms.
```