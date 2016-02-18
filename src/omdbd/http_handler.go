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
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UrlOptions struct {
	output_format int
	delete        bool
	showvalue     bool // used for list
	max           int  // used for list
}

func httpRequestHandler(w http.ResponseWriter, r *http.Request) (int, string) {
	g_info.IncreaseCounter("http.list")

	var err error
	var k *KeyInfo
	k, err = getKeyInfoFromUrlPath(r.URL.Path)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	var o *UrlOptions
	o, err = parseUrlOptions(r)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	// Is this put?
	paramV := r.Form["v"]
	if paramV != nil {
		return doPut(w, r, k, []byte(paramV[0]), o)
	}

	// Is thi list?
	if k.name[len(k.name)-1] == '/' {
		return doList(w, r, k, o)
	}

	// Do get
	return doGet(w, r, k, o)

	//root := strings.ToLower(getUrlValueStr(params, "root", ""))

	//depth := getUrlValueInt(params, "depth", 0)

	/*
		b, err := json.Marshal(results)
		if err != nil {
			g_info.IncreaseCounter("list.response.error")
			return http.StatusInternalServerError, err.Error()
		}

		g_info.IncreaseCounter("list.response.data")
	*/
}

func doPut(w http.ResponseWriter, r *http.Request, k *KeyInfo, v []byte, o *UrlOptions) (int, string) {
	timer := time.Now()

	if len(k.name) == 0 {
		return http.StatusBadRequest, "Empty keyname is not supported."
	}

	if err := dbf.Put(k, v); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	w.Write([]byte(time.Since(timer).String()))
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func doGet(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	timer := time.Now()

	v, err := dbf.Get(k)
	if err != nil {
		return http.StatusNotFound, err.Error()
	}

	if o.output_format == OUTPUT_FORMAT_URL {
		w.Write([]byte(url.QueryEscape(string(v))))
	} else {
		w.Write(v)
	}
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func doList(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	timer := time.Now()

	it := dbf.NewIterator(k.db)
	defer it.Close()

	start := string(k.name)
	start = start[0 : len(start)-1]
	g_debug.Println(start)
	if len(start) == 0 {
		it.SeekToFirst()
	} else {
		it.Seek([]byte(start))
	}
	for ; it.Valid(); it.Next() {
		w.Write(it.Key())
		if o.showvalue {
			w.Write([]byte("="))
			w.Write(it.Value())
		}
		w.Write([]byte("\n"))
	}

	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}

func parseUrlOptions(r *http.Request) (*UrlOptions, error) {
	options := UrlOptions{
		max: conf.Limit.ListSize,
	}

	o := r.Form["o"]
	if o != nil {
		for _, opt := range o {
			for _, v := range strings.Split(opt, ",") {
				if v == "text" {
					options.output_format = OUTPUT_FORMAT_TEXT
				} else if v == "url" {
					options.output_format = OUTPUT_FORMAT_URL
				} else if v == "delete" {
					options.delete = true
				} else if v == "showvalue" {
					options.showvalue = true
				} else if strings.HasPrefix(v, "max:") {
					var err error
					options.max, err = strconv.Atoi(v[4:len(v)])
					if err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("Unknown option. %s", v)
				}
			}
		}
	}

	return &options, nil
}
