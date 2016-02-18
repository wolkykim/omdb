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
	"runtime/debug"
	"sync"
	"testing"
)

func assert(t *testing.T, cond bool) {
	if !(cond) {
		debug.PrintStack()
		t.Fail()
	}
}

func TestProperty(t *testing.T) {
	info := New()

	info.SetProperty("property", "data")
	assert(t, "data" == info.GetProperty("property"))

	info.SetProperty("property", "data2")
	assert(t, "data2" == info.GetProperty("property"))
}

func TestCounter(t *testing.T) {
	info := New()

	assert(t, 0 == info.GetCounter("counter"))

	info.SetCounter("counter", 10)
	assert(t, 10 == info.GetCounter("counter"))

	info.IncreaseCounter("counter")
	assert(t, 11 == info.GetCounter("counter"))

	info.DecreaseCounter("counter")
	assert(t, 10 == info.GetCounter("counter"))
}

func TestConcurrency(t *testing.T) {
	const CONCURRENCY = 100
	const NUM_UPDATES = 10000

	info := New()
	var wgDone, wgTrigger sync.WaitGroup
	wgTrigger.Add(1)
	wgDone.Add(CONCURRENCY)

	for i := 0; i < CONCURRENCY; i++ {
		go func() {
			wgTrigger.Wait() // wait for the trigger
			for n := 0; n < NUM_UPDATES; n++ {
				info.IncreaseCounter("counter")
				info.DecreaseCounter("counter")
				info.UpdateCounter("counter", 1)
			}
			wgDone.Done()
		}()
	}

	wgTrigger.Done() // let all threads run
	wgDone.Wait() // wait all threads to finish
	assert(t, CONCURRENCY*NUM_UPDATES == info.GetCounter("counter"))
}
