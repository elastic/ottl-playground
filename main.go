// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
)

func main() {
	gzipHandler, err := gziphandler.NewGzipLevelHandler(gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":9090",
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           gzipHandler(http.FileServer(http.Dir("web/public"))),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
