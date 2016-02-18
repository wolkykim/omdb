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

package rotatelogger

import (
	"log"
	"os"
	"sync"
	"time"
)

type RotateLogger struct {
	lock             sync.Mutex
	log              *log.Logger
	out              *os.File
	filename         string
	rotationInterval int
	nextRotationTime time.Time
}

func New(filename string, rotationInterval int) *RotateLogger {
	l := &RotateLogger{
		filename:         filename,
		rotationInterval: rotationInterval,
		nextRotationTime: time.Now().Round(time.Duration(rotationInterval) * time.Second),
	}
	if err := l.checkAndRotate(); err != nil {
		return nil
	}
	return l
}

func (l *RotateLogger) Println(v ...interface{}) {
	l.checkAndRotate()
	l.lock.Lock()
	l.log.Println(v...)
	l.lock.Unlock()
}

func (l *RotateLogger) Printf(format string, v ...interface{}) {
	l.checkAndRotate()
	l.lock.Lock()
	l.log.Printf(format, v...)
	l.lock.Unlock()
}

func (l *RotateLogger) checkAndRotate() (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.out == nil || l.nextRotationTime.Before(time.Now()) {
		if l.out != nil {
			l.out.Close()
		}

		filename := l.nextRotationTime.Format(l.filename)

		// Set next rotation timing earlier to avoid frequent attempt when file open fails.
		l.nextRotationTime = l.nextRotationTime.Add(time.Duration(l.rotationInterval) * time.Second)
		out, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		l.out = out
		l.log = log.New(l.out, "", log.Ldate|log.Ltime)
	}
	return nil
}
