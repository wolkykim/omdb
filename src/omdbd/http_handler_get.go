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
	"time"
)

func doGet(w http.ResponseWriter, r *http.Request, k *KeyInfo, o *UrlOptions) (int, string) {
	g_info.IncreaseCounter("http.get")
	timer := time.Now()

	v, err := dbf.Get(k)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	if v == nil {
		return http.StatusNotFound, "No such key found."
	}

	var vm *ValueMeta
	vm, err = DecodeValue(v)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	w.Write(vm.Byte(o.encoding, k.path))
	return http.StatusOK, fmt.Sprint("runtime:", time.Since(timer))
}
