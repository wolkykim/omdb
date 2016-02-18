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
	"flag"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"omdbd/info"
	"log"
	"os"
	"time"
)

var g_info *info.Info
var g_debug *log.Logger
var g_exit bool
var g_exitcount int

func main() {
	var conf Config
	g_info = info.New()
	g_info.SetProperty("PRGNAME", PRGNAME)
	g_info.SetProperty("VERSION", VERSION)
	g_info.SetProperty("STARTED_AT", time.Now().Format(time.RFC3339))

	// Parse commend-line arguments
	if err := parseConfig(&conf); err != nil {
		log.Fatal(err)
	}

	// Init loggers
	if err := initLoggers(&conf); err != nil {
		log.Fatal(err)
	}

	// Launch HTTP listener
	log.Printf("Starting %s %s\n", PRGNAME, VERSION)
	if err := launchHttpListener(&conf); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(conf *Config) error {
	// Set command-line arguments
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "a filepath to the configuration file.")
	debug := flag.Bool("d", false, "Turn on debug output")
	flag.Parse()

	if err := gcfg.ReadFileInto(conf, *configFilepath); err != nil {
		return err
	}

	// Overwrite command-line parameters
	conf.Global.Debug = *debug
	return nil
}

// Initialize loggers
// We're using 2 loggers, default logger from log package for debug output
// and g_errlog for error logger.
func initLoggers(conf *Config) error {
	if conf.Global.Debug {
		// Enable debuglog
		g_debug = log.New(os.Stdout, "[DEBUG] ", log.Ltime)
		g_debug.Println("Entering debug mode.")

		// Set log output stream to stdout
		log.SetOutput(os.Stdout)
		log.SetFlags(log.Ltime)
		g_debug.Println("Redirecting log output to stdout.")
	} else {
		// Disable debuglog
		g_debug = log.New(ioutil.Discard, "DEBUG: ", log.Ltime)

		// Enable file logging
		if out, err := os.OpenFile(conf.Global.ErrorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		} else {
			log.SetOutput(out)
			log.SetFlags(log.Ldate | log.Ltime)
		}
	}
	return nil
}
