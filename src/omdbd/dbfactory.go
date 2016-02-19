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
	leveldb "github.com/jmhodges/levigo"
	"log"
	"regexp"
	"sync"
)

type DbFactory struct {
	dbdir      string
	autocreate bool
	cachesize  int

	opendbs map[string]*leveldb.DB
	ro      *leveldb.ReadOptions
	wo      *leveldb.WriteOptions
	mutex   sync.Mutex
}

func newDbFactory(config *Config) (*DbFactory, error) {
	dbfactory := &DbFactory{
		dbdir:      config.Database.Directory,
		autocreate: config.Database.AutoCreate,
		cachesize:  config.Database.CacheSize,

		opendbs: make(map[string]*leveldb.DB),
		ro:      leveldb.NewReadOptions(),
		wo:      leveldb.NewWriteOptions(),
	}
	g_debug.Printf("leveldb library version: %d.%d\n", leveldb.GetLevelDBMajorVersion(), leveldb.GetLevelDBMinorVersion())

	return dbfactory, nil
}

func (self *DbFactory) Close() {
	self.ro.Close()
	self.wo.Close()
	for dbname, db := range self.opendbs {
		log.Println("close database:", dbname)
		db.Close()
		g_info.IncreaseCounter("db.close")
	}
	log.Println("Closed all open databases.")
}

func (self *DbFactory) Open(dbname string) (*leveldb.DB, error) {
	if db := self.opendbs[dbname]; db != nil {
		return db, nil
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	if db := self.opendbs[dbname]; db != nil {
		return db, nil
	}

	// name check
	if isValidDbName(dbname) == false {
		return nil, fmt.Errorf("Invalid database name. %s", dbname)
	}

	cache := leveldb.NewLRUCache(self.cachesize)
	opts := leveldb.NewOptions()
	defer opts.Close()
	opts.SetCache(cache)
	opts.SetCreateIfMissing(self.autocreate)
	db, err := leveldb.Open(self.dbdir+"/"+dbname, opts)
	if err != nil {
		cache.Close()
		return nil, err
	}
	self.opendbs[dbname] = db
	log.Println("open database:", dbname)
	g_info.IncreaseCounter("db.open")
	return db, nil
}

func (self *DbFactory) NewIterator(dbname string) *leveldb.Iterator {
	db, err := self.Open(dbname)
	if err != nil {
		return nil
	}

	ro := leveldb.NewReadOptions()
	defer ro.Close()
	ro.SetFillCache(false)
	g_info.IncreaseCounter("db.iterator")
	return db.NewIterator(ro)
}

func (self *DbFactory) Get(k *KeyInfo) ([]byte, error) {
	db, err := self.Open(k.db)
	if err != nil {
		return nil, err
	}
	g_info.IncreaseCounter("db.get")
	return db.Get(self.ro, k.name)
}

func (self *DbFactory) Put(k *KeyInfo, v []byte) error {
	if conf.Default.MaxKeySize > 0 && len(k.name) > conf.Default.MaxKeySize {
		return fmt.Errorf("Key name is too long.")
	}
	if conf.Default.MaxValueSize > 0 && len(v) > conf.Default.MaxValueSize {
		return fmt.Errorf("Value size is too long.")
	}

	db, err := self.Open(k.db)
	if err != nil {
		return err
	}
	g_info.IncreaseCounter("db.put")
	return db.Put(self.wo, k.name, v)
}

func (self *DbFactory) Delete(k *KeyInfo) error {
	db, err := self.Open(k.db)
	if err != nil {
		return err
	}
	g_info.IncreaseCounter("db.del")
	return db.Delete(self.wo, k.name)
}

func isValidDbName(str string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_]*$").MatchString(str)
}
