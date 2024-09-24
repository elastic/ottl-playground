// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
)

func main() {
	wrapper, err := gziphandler.NewGzipLevelHandler(gzip.BestCompression)
	err = http.ListenAndServe(":9090", wrapper(http.FileServer(http.Dir("web/public"))))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
