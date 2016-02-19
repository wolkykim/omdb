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
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"net/url"
	"time"
)

const (
	VM_FLAG_ENVELOPE_VERSION = (0x01)    // version of serialized data format
	VM_FLAG_VERSIONNED       = (1 << 15) // indicate versioned data
)

type ValueMeta struct {
	flag   uint16 // first 4 bits are used for format numbering.
	Key    string `json:"k,omitempty"`
	Value  []byte `json:"v"`
	Ts     int64  `json:"ts,omitempty"`
	Expire int32  `json:"expire,omitempty"`
}

func (m *ValueMeta) Encode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(m.flag); err != nil {
		return nil, err
	}
	if err := encoder.Encode(m.Value); err != nil {
		return nil, err
	}
	if err := encoder.Encode(m.Ts); err != nil {
		return nil, err
	}
	if err := encoder.Encode(m.Expire); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (m *ValueMeta) Byte(encoding int, k string) []byte {
	if encoding == OUTPUT_ENCODING_URL {
		return []byte(url.QueryEscape(string(m.Value)))
	} else if encoding == OUTPUT_ENCODING_BASE64 {
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(m.Value)))
		base64.StdEncoding.Encode(encoded, m.Value)
		return encoded
	} else if encoding == OUTPUT_ENCODING_JSON {
		m.Key = k
		b, err := json.Marshal(m)
		if err != nil {
			return []byte("{\"error\":\"" + err.Error() + "\"")
		}
		return b
	} else {
		return m.Value
	}
}

func DecodeValue(buf []byte) (*ValueMeta, error) {
	var m ValueMeta
	decoder := gob.NewDecoder(bytes.NewBuffer(buf))
	if err := decoder.Decode(&m.flag); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&m.Value); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&m.Ts); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&m.Expire); err != nil {
		return nil, err
	}
	return &m, nil
}

func NewValue(v []byte) *ValueMeta {
	return &ValueMeta{
		flag:   VM_FLAG_ENVELOPE_VERSION,
		Value:  v,
		Ts:     time.Now().UnixNano(),
		Expire: 0,
	}
}
