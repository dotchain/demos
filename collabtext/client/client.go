// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/nw"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var argType = flag.String("type", "listx", "one of watch/list/listx/counter")
var argURL = flag.String("url", "http://:5000/api/", "server url")

func init() {
	flag.Parse()
}

func main() {
	client := &nw.Client{URL: *argURL}
	defer client.Close()

	tx := ops.TransformedWithCache(client, &sync.Map{})
	conn := ops.NewConnector(-1, nil, tx, rand.Float64)
	conn.Connect()

	switch *argType {
	case "watch":
		watch(client, conn)
	case "list":
		list(client)
	case "counter":
		count(client, conn)
	default:
		list(ops.Transformed(client))
	}

	conn.Disconnect()
}

func watch(client ops.Store, conn *ops.Connector) {
	val := text.StreamFromString("", false)

	streams.Connect(conn.Stream, val)

	val.Nextf("key", func() {
		count := 0
		for v, _ := val.Next(); v != nil; v, _ = val.Next() {
			val = v.(*text.Stream)
			count++
		}
		if count > 0 {
			log.Println("Value", val.Value())
		}
	})

	for {
		time.Sleep(5 * time.Second)
	}
}

func count(client ops.Store, conn *ops.Connector) {
	val := text.StreamFromString("", false).WithSessionID(ops.NewID())

	streams.Connect(conn.Stream, val)

	counter := 0
	for {
		for v, _ := val.Next(); v != nil; v, _ = val.Next() {
			val = v.(*text.Stream)
		}
		log.Println("Value", val.Value())
		counter++
		val = val.Paste(strconv.Itoa(counter))

		time.Sleep(5 * time.Second)
	}
}

func list(store ops.Store) {
	idx, limit := 0, 1000
	for idx%limit == 0 {
		ops, err := store.GetSince(getContext(10*time.Second), idx, limit)
		if err != nil {
			log.Fatal("Unexpected GetSince error", err)
		}

		if len(ops) == 0 {
			break
		}

		for _, op := range ops {
			log.Printf("log[%d] = %#v", idx, op)
			idx++
		}
	}
}

func getContext(duration time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	_ = cancel
	return ctx
}
