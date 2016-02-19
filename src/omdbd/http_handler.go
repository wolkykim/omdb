/******************************************************************************
 * OmDB
 *
 * Copyright (c) 2016 Seungyoung Kim.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *****************************************************************************/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UrlOptions struct {
	encoding  int
	showvalue bool // used for list
	max       int  // used for list
	filter    string
	remove    bool
}

func httpRequestHandler(w http.ResponseWriter, r *http.Request) (int, string) {
	var err error
	var k *KeyInfo
	k, err = getKeyInfoFromUrlPath(r.URL.Path)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	g_debug.Println(r.Method, r.RequestURI)

	var o *UrlOptions
	o, err = parseUrlOptions(r)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// Is this list?
	if k.name[len(k.name)-1] == '/' {
		return doList(w, r, k, o)
	}

	// Is this delete?
	if o.remove == true {
		return doDelete(w, r, k, o)
	}

	// Is this put?
	if r.Method == "PUT" {
		if r.ContentLength <= 0 {
			return http.StatusBadRequest, "No data to set."
		}
		return doPut(w, r, k, getHttpBody(r), o)
	} else if r.FormValue("v") != "" {
		return doPut(w, r, k, []byte(r.FormValue("v")), o)
	}

	// Do get
	return doGet(w, r, k, o)
}

func doPut(w http.ResponseWriter, r *http.Request, k *KeyInfo, v []byte, o *UrlOptions) (int, string) {
	g_info.IncreaseCounter("http.put")
	timer := time.Now()

	if len(k.name) == 0 {
		return http.StatusBadRequest, "Key name is empty."
	}

	// Encode
	vm := NewValue(v)
	vb, err := vm.Encode()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	// Store
	if err := dbf.Put(k, vb); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func doGet(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	g_info.IncreaseCounter("http.get")
	timer := time.Now()

	v, err := dbf.Get(k)
	if err != nil {
		g_info.IncreaseCounter("http.get.500")
		return http.StatusInternalServerError, err.Error()
	}
	if v == nil {
		g_info.IncreaseCounter("http.get.404")
		return http.StatusNotFound, "No such key found."
	}

	var vm ValueMeta
	if err = DecodeValue(&vm, v); err != nil {
		g_info.IncreaseCounter("http.get.500")
		return http.StatusInternalServerError, err.Error()
	}

	g_info.IncreaseCounter("http.get.200")
	w.Write(vm.Byte(o.encoding, k.path))
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func doDelete(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	g_info.IncreaseCounter("http.delete")
	timer := time.Now()

	if len(k.name) == 0 {
		return http.StatusBadRequest, "Key name is empty."
	}

	// Delete
	if err := dbf.Delete(k); err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func doList(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	g_info.IncreaseCounter("http.list")
	timer := time.Now()

	it := dbf.NewIterator(k.db)
	defer it.Close()

	start := string(k.name)
	start = start[0 : len(start)-1]
	if len(start) == 0 {
		it.SeekToFirst()
	} else {
		it.Seek([]byte(start))
	}

	b := new(bytes.Buffer)
	var vms []ValueMeta
	for n := 0; it.Valid(); it.Next() {
		key := genKeyPath(k.db, string(it.Key()))
		if o.max > 0 && n >= o.max {
			w.Header().Add(HTTP_HEADER_X_TRUNCATED, "1")
			w.Header().Add(HTTP_HEADER_X_NEXT, urlencode(key+"/"))
			break
		}

		// Filter
		if o.filter != "" {
			matched, err := regexp.MatchString(o.filter, key)
			if err != nil {
				return http.StatusBadRequest, err.Error()
			}
			if matched == false {
				continue
			}
		}

		if o.encoding == OUTPUT_ENCODING_JSON {
			var vm ValueMeta
			if err := DecodeValue(&vm, it.Value()); err != nil {
				return http.StatusInternalServerError, err.Error()
			}
			vm.Key = key
			if o.showvalue == false {
				vm.Value = nil
			}
			vms = append(vms, vm)
		} else {
			b.Write(encodeKey(o.encoding, []byte(key)))
			if o.showvalue {
				var vm ValueMeta
				if err := DecodeValue(&vm, it.Value()); err != nil {
					return http.StatusInternalServerError, err.Error()
				}
				b.Write([]byte("="))
				b.Write(vm.Byte(o.encoding, ""))
			}
			b.Write([]byte("\n"))
		}
		n = n + 1
	}

	if len(vms) > 0 {
		bb, err := json.MarshalIndent(vms, "", "\t")
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		w.Write(bb)
	} else {
		w.Write(b.Bytes())
	}
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func parseUrlOptions(r *http.Request) (*UrlOptions, error) {
	options := UrlOptions{
		max:      conf.Default.ListSize,
		encoding: conf.Default.OutputEncoding,
	}

	o := r.Form["o"]
	if o != nil {
		for _, opt := range o {
			for _, v := range strings.Split(opt, ",") {
				if v == "binary" {
					options.encoding = OUTPUT_ENCODING_BINARY
				} else if v == "url" {
					options.encoding = OUTPUT_ENCODING_URL
				} else if v == "base64" {
					options.encoding = OUTPUT_ENCODING_BASE64
				} else if v == "json" {
					options.encoding = OUTPUT_ENCODING_JSON
				} else if v == "delete" {
					options.remove = true
				} else if v == "showvalue" {
					options.showvalue = true
				} else if strings.HasPrefix(v, "max:") {
					var err error
					options.max, err = strconv.Atoi(v[len("max:"):len(v)])
					if err != nil {
						return nil, err
					}
				} else if strings.HasPrefix(v, "filter:") {
					options.filter = v[len("filter:"):len(v)]
				} else {
					return nil, fmt.Errorf("Unknown option. %s", v)
				}
			}
		}
	}

	return &options, nil
}
