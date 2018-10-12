// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"github.com/dotchain/dot/x/nw"
	"log"
	"net/http"
	"time"
)

func main() {
	store := nw.MemPoller(nw.MemStore(nil))
	defer store.Close()
	handler := &nw.Handler{Store: store}
	h := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, req)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
	})
	mux.HandleFunc("/api/", h)
	mux.HandleFunc("/api", h)
	srv := &http.Server{
		Addr:           ":5000",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 10000,
	}
	log.Fatal(srv.ListenAndServe())
}
