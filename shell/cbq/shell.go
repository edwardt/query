//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"crypto/tls"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tiServer = flag.String("engine", "http://localhost:8093/", "URL to cbq-engine")

func main() {
	flag.Parse()
	if strings.HasSuffix(*tiServer, "/") == false {
		*tiServer = *tiServer + "/"
	}
	HandleInteractiveMode(*tiServer, filepath.Base(os.Args[0]))
}

func execute_internal(tiServer, line string, w *os.File) error {

	url := tiServer + "query"
	tr := &http.Transport{}
	if strings.HasPrefix(url, "https") {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "text/plain", strings.NewReader(line))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
	w.WriteString("\n")
	w.Sync()

	return nil
}
