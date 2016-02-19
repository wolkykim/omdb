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
	"time"
)

const (
	HTTP_HEADER_X_TRUNCATED = "X-Omdb-Truncated"
	HTTP_HEADER_X_NEXT      = "X-Omdb-Next"
)

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
	scanCount := 0;
	listCount := 0;
	for ; it.Valid(); it.Next() {
		scanCount = scanCount + 1

		key, _ := parseKeyInfo(genKeyPath(k.db, string(it.Key())))
		if (o.max > 0 && listCount >= o.max) || (o.maxscan > 0 && scanCount >= o.maxscan) {
			w.Header().Add(HTTP_HEADER_X_TRUNCATED, "1")
			w.Header().Add(HTTP_HEADER_X_NEXT, urlencode(key.path+"/"))
			break
		}

		// Filter versioned keys
		if o.showversion == false && key.IsVersion() == true {
			continue
		}

		// Apply user filter
		if o.filter != "" {
			matched, err := regexp.MatchString(o.filter, key.path)
			if err != nil {
				return http.StatusBadRequest, err.Error()
			}
			if matched == false {
				continue
			}
		}

		listCount = listCount + 1

		if o.remove == true {
			deleteKey(key)
			continue
		}

		if o.encoding == OUTPUT_ENCODING_JSON {
			vm, err := DecodeValue(it.Value())
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}
			vm.Key = key.path
			if o.showvalue == false {
				vm.Value = nil
			}
			vms = append(vms, *vm)
		} else {
			b.Write(key.Byte(o.encoding, o.html))
			if o.showvalue {
				vm, err := DecodeValue(it.Value())
				if err != nil {
					return http.StatusInternalServerError, err.Error()
				}
				b.Write([]byte("="))
				b.Write(vm.Byte(o.encoding, ""))
			}
			b.Write([]byte("\n"))
		}
	}

	if len(vms) > 0 {
		bb, err := json.MarshalIndent(vms, "", "\t")
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		w.Write(bb)
	} else if b.Len() > 0 {
		w.Write(b.Bytes())
	}
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}
