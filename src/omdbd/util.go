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
	"strings"
)

func isWebBrowser(r *http.Request) bool {
	agent := r.Header.Get("User-Agent")
	if len(agent) == 0 {
		return false
	}
	if strings.HasPrefix(agent, "Mozilla/5.0 (") {
		if strings.Contains(agent, " Firefox/") == true ||
			strings.Contains(agent, " Chrome/") == true ||
			strings.Contains(agent, "MSIE ") == true {
			return true
		}
	}
	return false
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
