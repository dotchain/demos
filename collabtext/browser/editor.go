// +js
// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/dotchain/dot/ops"
	"github.com/dotchain/dot/streams"
	"github.com/dotchain/dot/streams/text"
	"github.com/dotchain/dot/x/idgen"
	"github.com/dotchain/dot/x/nw"
	"github.com/dotchain/dot/x/types"
	"github.com/gopherjs/gopherjs/js"
	"log"
	"sync"
	"time"
)

func NewEditable(url string, done func(map[string]interface{}), refresh func(text string, start, end int)) {
	go func() {
		client := &nw.Client{URL: url}
		cache := &sync.Map{}
		tx := ops.TransformedWithCache(client, cache)
		stream := streams.New()
		val := text.StreamFromString("", false)
		ch := make(chan func(), 1000)
		b := &streams.Branch{stream, val.WithoutOwnCursor()}
		sync := ops.NewSync(tx, -1, stream, idgen.New)

		defer client.Close()
		defer sync.Close()

		e := &editable{b, val, sync, refresh}
		done(e.jsObject(ch))

		go e.poll(ch, client)
		for fn := range ch {
			fn()
		}
	}()
}

type editable struct {
	b       *streams.Branch
	val     *text.Stream
	sync    *ops.Sync
	refresh func(text string, start, end int)
}

func (e *editable) update() {
	b, val, refresh := e.b, e.val, e.refresh

	b.Merge()
	for _, v := val.Next(); v != nil; _, v = val.Next() {
		val = v.(*text.Stream)
	}
	e.val = val
	start, _ := val.E.Start()
	end, _ := val.E.End()
	text := types.S16(val.E.Text)
	start, end = text.ToUTF16(start), text.ToUTF16(end)
	log.Println("Value", val.E.Text, start, end)
	refresh(val.E.Text, start, end)
}

func (e *editable) fetch(ch chan func()) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	fetched, err := e.sync.Prefetch(ctx, 1000)
	if err != nil || !fetched {
		if err != nil && ctx.Err() == nil {
			log.Println("Unexpected GetSince() error", err)
		}
		return
	}

	ch <- func() {
		e.sync.ApplyPrefetched()
		e.update()
	}
}

func (e *editable) poll(ch chan func(), store ops.Store) {
	for {
		version := e.sync.Version()
		e.fetch(ch)
		if e.sync.Version() != version {
			continue
		}
		duration := time.Second * 30
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		err := store.Poll(ctx, version+1)
		if err != nil && ctx.Err() == nil {
			log.Println("Unexpected poll error", err)
		}
		cancel()
	}
}

func (e *editable) jsObject(ch chan func()) map[string]interface{} {
	return map[string]interface{}{
		"Insert": func(s string) {
			e.val.Insert(s)
			ch <- e.update
		},
		"Delete": func() {
			e.val.Delete()
			ch <- e.update
		},
	}
}

func main() {
	js.Global.Set("NewEditable", NewEditable)
}
