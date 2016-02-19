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
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	OUTPUT_ENCODING_BINARY = 0
	OUTPUT_ENCODING_URL    = 1
	OUTPUT_ENCODING_BASE64 = 2
	OUTPUT_ENCODING_JSON   = 3
)

type UrlOptions struct {
	encoding int
	html     bool
	remove   bool

	// used for list
	showvalue   bool
	showversion bool
	max         int
	maxscan     int
	filter      string
}

func httpRequestHandler(w http.ResponseWriter, r *http.Request) (int, string) {
	var err error
	var k *KeyInfo
	k, err = parseKeyInfo(r.URL.Path)
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

func parseUrlOptions(r *http.Request) (*UrlOptions, error) {
	options := UrlOptions{
		encoding: OUTPUT_ENCODING_BINARY,
	}

	o := []string{conf.Global.QueryOption}
	o = append(o, r.Form["o"]...)
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
				} else if v == "html" {
					options.html = true
				} else if v == "delete" {
					options.remove = true
				} else if v == "showvalue" {
					options.showvalue = true
				} else if v == "showversion" {
					options.showversion = true
				} else if strings.HasPrefix(v, "max:") {
					var err error
					options.max, err = strconv.Atoi(v[len("max:"):len(v)])
					if err != nil {
						return nil, err
					}
				} else if strings.HasPrefix(v, "maxscan:") {
					var err error
					options.maxscan, err = strconv.Atoi(v[len("maxscan:"):len(v)])
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
