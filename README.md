# demos
demos for dot/chain project

[![Status](https://travis-ci.com/dotchain/demos.svg?branch=master)](https://travis-ci.com/dotchain/demos?branch=master)
[![GoDoc](https://godoc.org/github.com/dotchain/demos?status.svg)](https://godoc.org/github.com/dotchain/demos)
[![codecov](https://codecov.io/gh/dotchain/demos/branch/master/graph/badge.svg)](https://codecov.io/gh/dotchain/demos)
[![GoReportCard](https://goreportcard.com/badge/github.com/dotchain/demos)](https://goreportcard.com/report/github.com/dotchain/demos)


## Page

The demos can be available online [here](https://dotchain.github.io/demos/)

## CollabText

The collab test can be run locally. Multiple browser and local clients can be used.

Start the server:

```sh
go run collabtext/server/server.go
```

Start a client watcher:

```sh
go run collabtext/client/client.go -type watch
```

Start a client counter:

```sh
go run collabtext/client/client.go -type counter
```

Multiple counters can be started illustrating how synchronization works.

Start a local http server:

```sh
$> cd collabtext/browser
$> ln -s  ../../../../ github.com
$> python -m SimpleHTTPServer 5001
$> open http://localhost:5001
```

Note that the browser and the counter work on the same "document"
