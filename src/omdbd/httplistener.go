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
	"errors"
	"fmt"
	"log"
	"net/http"
	"omdbd/rotatelogger"
	"os"
	"os/signal"
	"syscall"
)

type handler func(w http.ResponseWriter, r *http.Request)
type reqhandler func(w http.ResponseWriter, r *http.Request) (int, string)

var accesslog *rotatelogger.RotateLogger
var dbf *DbFactory
var conf *Config

func launchHttpListener(config *Config) error {
	// Check and get new db connector
	conf = config
	var err error
	if dbf, err = newDbFactory(conf); err != nil {
		return err
	}

	// Open access log
	accesslog = rotatelogger.New(conf.Httplistener.AccessLog, conf.Httplistener.LogRotate)
	if accesslog == nil {
		return fmt.Errorf("Failed to open %s", conf.Httplistener.AccessLog)
	}

	// Create a pidfile
	if err = createPidfile(conf.Global.PidFile, os.Getpid()); err != nil {
		return err
	}
	defer removePidfile(conf.Global.PidFile)

	// Start to handle signals
	initSignals()

	// Set up and launch HTTP listener
	http.HandleFunc(conf.Global.StatusUrl, superHandler(httpStatusHandler))
	http.HandleFunc("/", validator(superHandler(httpRequestHandler)))

	log.Printf("%s started - listening on %s://%s\n", PRGNAME, conf.Httplistener.Protocol, conf.Httplistener.Addr)
	if conf.Httplistener.Protocol == "http" {
		return http.ListenAndServe(conf.Httplistener.Addr, nil)
	} else if conf.Httplistener.Protocol == "https" {
		return http.ListenAndServeTLS(conf.Httplistener.Addr,
			conf.Httplistener.CertFile, conf.Httplistener.KeyFile, nil)
	}
	return errors.New("Not supported protocol." + conf.Httplistener.Protocol)
}

// This is an wrapper to validate requests.
func validator(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement request validation code.
		var valid = true // bypass for now
		if valid == true {
			h(w, r)
			return
		}
		http.Error(w, "Permission denied.", http.StatusForbidden)
	}
}

// This is an wrapper to pre/post processing requests for common
func superHandler(h reqhandler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update common headers before user handler, in case, http.Error() is called,
		// it flushes out right away, so updating headers after user handler may not work.
		w.Header().Add("Server", fmt.Sprintf("%s (v%s)", PRGNAME, VERSION))

		var status int
		var message string
		if g_exit == false {
			r.ParseForm()
			status, message = h(w, r)
		} else {
			status = http.StatusServiceUnavailable
			message = "The server is shutting down."
		}

		if (status / 100) != 2 {
			http.Error(w, message, status)
		}

		accesslog.Printf("%s \"%s %s\" %d \"%s\"\n", r.RemoteAddr, r.Method, r.RequestURI, status, message)
		g_info.IncreaseCounter(fmt.Sprintf("http.rescode.%d", status))
	}
}

func initSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		g_exit = true
		log.Println("Signal received -", sig)
		if dbf != nil {
			dbf.Close()
		}
		log.Printf("%s terminated.\n", PRGNAME)
		removePidfile(conf.Global.PidFile)
		os.Exit(0)
	}()
}
