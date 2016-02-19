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
	"regexp"
	"strings"
)

const (
	VERSIONING_DELIMITER = "#V"
)

type KeyInfo struct {
	path string
	db   string
	name []byte
}

func (k *KeyInfo) Byte(encoding int, html bool) []byte {
	var key string
	if encoding == OUTPUT_ENCODING_BINARY {
		key = k.path
	} else {
		key = urlencode(string(k.path))
	}

	if html == true {
		key = fmt.Sprintf("<a href=\"%s\">%s</a><br>", urlencode(string(k.path)), key)
	}

	return []byte(key)
}

func (k *KeyInfo) IsVersion() bool {
	matched, err := regexp.MatchString(".+"+VERSIONING_DELIMITER+"[0-9]{19}$", k.path)
	if err != nil {
		return false
	}
	return matched
}

func parseKeyInfo(urlpath string) (*KeyInfo, error) {
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

func genKeyPath(db string, name string) string {
	return "/" + db + "/" + name
}

func genVersionKeyInfo(k *KeyInfo, ts int64) (*KeyInfo, error) {
	verno := MAX_INT64 - ts
	return parseKeyInfo(genKeyPath(k.db, fmt.Sprintf("%s%s%019d", string(k.name), VERSIONING_DELIMITER, verno)))
}
