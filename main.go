// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	relativeWebPublicDir = "web/public"
	defaultAddr          = ":8080"
)

func main() {
	listenAddress, ok := os.LookupEnv("ADDR")
	if !ok {
		listenAddress = defaultAddr
	}

	mux := http.NewServeMux()
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
