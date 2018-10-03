// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/refs"
	"github.com/dotchain/dot/x/idgen"
	"github.com/dotchain/dot/x/nw"
	"github.com/dotchain/dot/x/types"
	"log"
	"strconv"
	"time"
)

var argType = flag.String("type", "listx", "one of watch/list/listx/counter")
var argUrl = flag.String("url", "http://:8183", "server url")

func init() {
	flag.Parse()
}

func main() {
	client := &nw.Client{URL: *argUrl}
	defer client.Close()

	tx := ops.Transformed(client)
	stream := changes.NewStream()

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

func watch(client ops.Store, sync *ops.Sync, stream changes.Stream) {
	var value changes.Value
	var version int

	value = types.S8("")
	for {
		ctx := getContext(30 * time.Second)
		if err := sync.Fetch(ctx, 1000); err != nil {
			log.Fatal("Unexpected fetch error", err)
		}

		for cx, sx := stream.Next(); sx != nil; cx, sx = stream.Next() {
			stream = sx
			value = value.Apply(cx)
			version++
		}
		log.Println("Value", value)

		ctx = getContext(30 * time.Second)
		if err := client.Poll(ctx, version); err != nil && ctx.Err() == nil {
			log.Fatal("Unexpected poll error", err)
		}
	}
}

func count(client ops.Store, sync *ops.Sync, stream changes.Stream) {
	var ref refs.Ref
	var value changes.Value
	var version, counter int

	value, ref = types.S8(""), refs.Caret{nil, 0}
	for {
		ctx := getContext(30 * time.Second)
		if err := sync.Fetch(ctx, 100); err != nil {
			log.Fatal("Unexpected fetch error", err)
		}

		for cx, sx := stream.Next(); sx != nil; cx, sx = stream.Next() {
			stream = sx
			value = value.Apply(cx)
			ref, _ = ref.Merge(cx)
			version++
		}
		log.Println("Value", value)

		offset := 0
		if caret, ok := ref.(refs.Caret); ok {
			offset = caret.Index
		}

		before := types.S8(strconv.Itoa(counter))
		if counter == 0 {
			before = types.S8("")
		} else if value.Slice(offset, before.Count()) != before {
			log.Fatal("Unexpected offset", offset, value)
		}

		counter++
		after := types.S8(strconv.Itoa(counter))

		ch := changes.Splice{offset, before, after}
		stream = stream.Append(ch)
		value = value.Apply(ch)
		version++

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