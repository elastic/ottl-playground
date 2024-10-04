// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	relativeWebPublicDir = "web/public"
	defaultAddr          = ":8080"
	webAssemblyFileName  = "ottlplayground.wasm"
)

func main() {
	listenAddress, ok := os.LookupEnv("ADDR")
	if !ok {
		listenAddress = defaultAddr
	}

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s", webAssemblyFileName), newWebAssemblyHandler())
	mux.HandleFunc("/", http.FileServer(http.Dir(webPublicDir())).ServeHTTP)

	server := &http.Server{
		Addr:              listenAddress,
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           mux,
	}

	log.Println("Listening on ", listenAddress)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func newWebAssemblyHandler() http.Handler {
	wasmFile := filepath.Join(webPublicDir(), webAssemblyFileName)
	file, err := os.Open(wasmFile)
	if err != nil {
		log.Fatal(err)
	}

	detectBuff := make([]byte, 512)
	_, err = file.Read(detectBuff)
	if err != nil {
		log.Fatal(err)
	}

	_ = file.Close()

	contentType := http.DetectContentType(detectBuff)
	if contentType == "application/x-gzip" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/wasm")
			w.Header().Set("Content-Encoding", "gzip")
			http.ServeFile(w, r, wasmFile)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, wasmFile)
	})
}

func webPublicDir() string {
	if _, err := os.Stat(relativeWebPublicDir); err == nil {
		return relativeWebPublicDir
	}
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(filepath.Dir(executable), relativeWebPublicDir)
}
