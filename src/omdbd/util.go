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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func getKeyInfoFromUrlPath(urlpath string) (*KeyInfo, error) {
	path := strings.TrimSpace(urlpath)
	tokens := strings.SplitN(path, "/", 3)
	if len(tokens) != 3 {
		return nil, fmt.Errorf("Invalid urlpath. %s", urlpath)
	}

	name := tokens[2]
	if len(name) == 0 {
		name = "/"
	}

	return &KeyInfo{
		path: path,
		db:   tokens[1],
		name: []byte(name),
	}, nil
}

func getHttpBody(r *http.Request) []byte {
	if r.ContentLength == 0 {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	return body
}

func urlencode(s string) string {
	return strings.Replace(url.QueryEscape(s), "%2F", "/", -1)
}

func encodeKey(encoding int, k []byte) []byte {
	if encoding == OUTPUT_ENCODING_BINARY {
		return k
	} else {
		return []byte(urlencode(string(k)))
	}
}

func genKeyPath(db string, keypath string) string {
	return "/" + db + "/" + keypath
}

func floatToString(f float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func getUrlValueStr(urlValues url.Values, name string, defvalue string) string {
	params := urlValues[name]
	if params == nil {
		return defvalue
	}
	return params[0]
}

func getUrlValueInt(urlValues url.Values, name string, defvalue int) int {
	num, err := strconv.Atoi(getUrlValueStr(urlValues, name, strconv.Itoa(defvalue)))
	if err != nil {
		return defvalue
	}
	return num
}

func createPidfile(pidpath string, pid int) error {
	fp, err := os.OpenFile(pidpath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err = fmt.Fprint(fp, pid); err != nil {
		return err
	}
	return nil
}

func removePidfile(pidpath string) error {
	return os.Remove(pidpath)
}
