/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Elasticsearch B.V. licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
