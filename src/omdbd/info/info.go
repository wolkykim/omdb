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

package info

import (
	"strconv"
	"sync"
)

type Info struct {
	properties map[string]string
	counters   map[string]int
	mutex      sync.Mutex
}

func New() *Info {
	var info Info
	info.properties = make(map[string]string)
	info.counters = make(map[string]int)
	return &info
}

func (info *Info) SetProperty(key string, value string) string {
	info.mutex.Lock()
	info.properties[key] = value
	info.mutex.Unlock()
	return value
}

func (info *Info) GetProperty(key string) string {
	info.mutex.Lock()
	value := info.properties[key]
	info.mutex.Unlock()
	return value
}

func (info *Info) SetCounter(key string, value int) int {
	info.mutex.Lock()
	info.counters[key] = value
	info.mutex.Unlock()
	return value
}

func (info *Info) GetCounter(key string) int {
	info.mutex.Lock()
	value := info.counters[key]
	info.mutex.Unlock()
	return value
}

func (info *Info) UpdateCounter(key string, inc int) int {
	info.mutex.Lock()
	value := info.counters[key] + inc
	info.counters[key] = value
	info.mutex.Unlock()
	return value
}

func (info *Info) IncreaseCounter(key string) int {
	return info.UpdateCounter(key, 1)
}

func (info *Info) DecreaseCounter(key string) int {
	return info.UpdateCounter(key, -1)
}

// Returns a map that contains both properties and counters.
// The counter values will be stored in string value.
func (info *Info) ToMap() map[string]string {
	newmap := make(map[string]string)

	info.mutex.Lock()
	for k, v := range info.properties {
		newmap[k] = v
	}
	for k, v := range info.counters {
		newmap[k] = strconv.Itoa(v)
	}
	info.mutex.Unlock()

	return newmap
}
