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
	"github.com/dotchain/dot/x/idgen"
	"github.com/dotchain/dot/x/nw"
	"log"
	"strconv"
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

	tx := ops.Transformed(client)
	stream := streams.New()

	sync := ops.NewSync(tx, -1, stream, idgen.New)
	defer sync.Close()

	switch *argType {
	case "watch":
		watch(client, sync, stream)
	case "list":
		list(client)
	case "counter":
		count(client, sync, stream)
	default:
		list(ops.Transformed(client))
	}
}

func watch(client ops.Store, sync *ops.Sync, stream streams.Stream) {
	val := text.StreamFromString("", false)

	b := streams.Branch{stream, val.WithoutOwnCursor()}
	b.Connect()

	version := 0
	for {
		ctx := getContext(10 * time.Second)
		if err := sync.Fetch(ctx, 1000); err != nil {
			log.Fatal("Unexpected fetch error", err)
		}

		for _, v := val.Next(); v != nil; _, v = val.Next() {
			val = v.(*text.Stream)
			version++
		}
		log.Println("Value", val.E.Text)

		ctx = getContext(10 * time.Second)
		if err := client.Poll(ctx, version); err != nil && ctx.Err() == nil {
			log.Fatal("Unexpected poll error", err)
		}
	}
}

func count(client ops.Store, sync *ops.Sync, stream streams.Stream) {
	val := text.StreamFromString("", false)

	b := streams.Branch{stream, val.WithoutOwnCursor()}
	b.Connect()

	counter := 0
	for {
		ctx := getContext(30 * time.Second)
		if err := sync.Fetch(ctx, 100); err != nil {
			log.Fatal("Unexpected fetch error", err)
		}

		for _, v := val.Next(); v != nil; _, v = val.Next() {
			val = v.(*text.Stream)
		}
		log.Println("Value", val.E.Text)

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
