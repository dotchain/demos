// +js
// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/nw"
	"github.com/gopherjs/gopherjs/js"
	"log"
	"math/rand"
	"sync"
)

// NewEditable asynchronously returns an interface which the browser
// can use to send events to. The refresh method is called (also
// asynchronously) in response to every update of local state.
func NewEditable(url string, done func(*js.Object), refresh func(*js.Object)) {
	go func() {
		client := &nw.Client{URL: url}
		cache := &sync.Map{}
		tx := ops.TransformedWithCache(client, cache)
		conn := ops.NewConnector(-1, nil, tx, rand.Float64)
		val := text.StreamFromString("", false)
		val.S = conn.Async.Wrap(val.S)
		wo := val.WithoutOwnCursor()
		streams.Connect(conn.Stream, wo)

		// defer client.Close()
		// defer conn.Disconnect()

		val.Nextf("key", func() {
			for v, _ := val.Next(); v != nil; v, _ = val.Next() {
				val = v.(*text.Stream)
			}
			start, _ := val.Start(true)
			end, _ := val.End(true)
			log.Println("Value", val.Value(), start, end)
			refresh(js.MakeWrapper(val))
		})

		done(js.MakeWrapper(val))
		conn.Connect()
	}()
}

func main() {
	js.Global.Set("NewEditable", NewEditable)
}
