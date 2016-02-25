About OmDB
==========

* Persistent Key/Value Store.
* Simple in distribution – One binary executable.
* No special client but just any HTTP/HTTPS clients, even web browsers..
* Sorted Key Listing and Iteration.
* Versioning (Yeah~)
* Written in Go lang with leveldb storage backend.

Build & Install
===============
```
$ git clone https://github.com/wolkykim/omdb.git
$ cd omdb
$ make build-deps
$ sudo make install-deps
$ make

$ sudo make install
$ cd /usr/local/omdb;
$ sudo chown RUN_USER:RUN_GROUP db logs
$ sudo cp etc/omdbd.conf.example etc/omdbd.conf
$ bin/omdbd -c etc/omdbd.conf -d
```

APIs
====

## Key Insertion & Retrieval

OmDB takes PUT, GET and POST methods for key insertion as below.

List all keys in sorted order that are prefixed by a given key or within the given range.
```
$ cat test.txt
Hello World!
This is OmDB

$ curl http://localhost:8081/testdb/test.txt -T test.txt

$ curl http://localhost:8081/testdb/test.txt
Hello World!
This is OmDB
```

Here's another way of insertion.

```
$ curl http://localhost:8081/testdb/memo --data "v=I'm still hungry."
$ curl http://localhost:8081/testdb/memo                             
I'm still hungry.
```

### Output encodings

It supports url, base64 and json encodings.

```
$ curl http://localhost:8081/testdb/test.txt --data "o=url"
Hello+World%21%0AThis+is+OmDB%0A

$ curl http://localhost:8081/testdb/test.txt --data "o=base64"
SGVsbG8gV29ybGQhClRoaXMgaXMgT21EQgo=

$ curl http://localhost:8081/testdb/test.txt --data "o=json"
{
	"k": "/testdb/test.txt",
	"v": "SGVsbG8gV29ybGQhClRoaXMgaXMgT21EQgo=",
	"ts": 1456371872563194746
}

Try 'html' option when working with web browsers.
```

## Key Listing and Regular-expression Search

When given key has tailing slash, it lists keys starting from the given key or the very next key if no matching key found.
Keys are listed in sorted manner.

```
$ curl http://localhost:8081/testdb/                                       
/testdb/memo
/testdb/mypic
/testdb/notes/tobuy
/testdb/notes/todo
/testdb/test.txt

$ curl http://localhost:8081/testdb/notes/tod/
/testdb/notes/todo
/testdb/test.txt
```

### Show values - 'showvalue' option

```
$ curl http://localhost:8081/testdb/notes/tod/ --data "o=showvalue,url"
/testdb/notes/todo=I%27m+still+hungry.
/testdb/test.txt=Hello+World%21%0AThis+is+OmDB%0A

In Json encoding, the values get Base64 encoded.

$ curl http://localhost:8081/testdb/notes/tod/ --data "o=showvalue,json"
[
	{
		"k": "/testdb/notes/todo",
		"v": "SSdtIHN0aWxsIGh1bmdyeS4=",
		"ts": 1456372221395798559
	},
	{
		"k": "/testdb/test.txt",
		"v": "SGVsbG8gV29ybGQhClRoaXMgaXMgT21EQgo=",
		"ts": 1456371872563194746
	}
]
```

### Search - 'filter' and 'maxscan' options

```
$ curl http://localhost:8081/testdb/notes/tod/ --data "o=showvalue,url,filter:.*no.*"
/testdb/notes/todo=I%27m+still+hungry.
```

'maxscan' options can be used to limit internal scan range. This is useful to control the load when filter option is given and there are lesser number of matching keys thax 'max' size, it'll limit number of key iteration to 'maxscan' size otherwise it'll continue the scan till the end.

### Limiting list size - 'max' option

'max' option specifies the number of keys in the return. If there are more keys that this, it will attach 'X-Omdb-Truncated' and 'X-Omdb-Next' headers for next iteration. By default, it's set to 1000.

```
$ curl http://localhost:8081/testdb/ --data "o=max:2"
(Response Header) X-Omdb-Truncated: 1
(Response Header) X-Omdb-Next: /testdb/mypic%23V7766999835404987155/
/testdb/memo
/testdb/mypic

$ curl http://localhost:8081/testdb/mypic%23V7766999835404987155/ --data "o=max:2"
(Response Header) X-Omdb-Next: /testdb/notes/todo%23V7766999815458977248/
(Response Header) X-Omdb-Truncated: 1
/testdb/notes/tobuy
/testdb/notes/todo

$ curl http://localhost:8081/testdb/notes/todo%23V7766999815458977248/ --data "o=max:2"   
/testdb/test.txt
```

## Versioning

Versioning is a configurable optional feature.

```
$ curl http://localhost:8081/testdb/test.txt/ --data "o=max:1"
/testdb/test.txt

$ curl http://localhost:8081/testdb/test.txt/ --data "o=max:3,showversion,json"
[
	{
		"k": "/testdb/test.txt",
		"ts": 1456373530935729691
	},
	{
		"k": "/testdb/test.txt#V7766998505919046116",
		"ts": 1456373530935729691
	},
	{
		"k": "/testdb/test.txt#V7767000164291581061",
		"ts": 1456371872563194746
	}
]

$ curl http://localhost:8081/testdb/test.txt#V7766998505919046116
hello~
```
'#V(19 digit rev number)' is a special key postfix. As you see rev number decreases while timestamp(ts) increases. So lastest revision always located on the top right below to the key and the oldest one gets located on the bottom in the list.

## Deletion

```
$ curl http://localhost:8081/testdb/test.txt --data "o=delete"
$ curl http://localhost:8081/testdb/test.txt                  
No such key found.

$ curl http://localhost:8081/testdb/ --data "o=max:2"
/testdb/notes/tobuy
/testdb/notes/todo
$ curl http://localhost:8081/testdb/ --data "o=max:2,delete"
```

## Internal status page

```
$ curl http://localhost:8081/status                              
PRGNAME : omdbd
STARTED : 2016-02-24T19:44:00-08:00
VERSION : 1.0.0
db.get : 11
db.iterator : 34
db.open : 1
db.put : 12
http.get : 11
http.list : 36
http.put : 6
http.rescode.200 : 48
http.rescode.400 : 3
http.rescode.404 : 3
http.status : 1
```

## Available options

Option format is "o=option,option:val,..."

* binary - no encoding
* url - encode value with url encoding
* base64 - encode value with base74 encoding. keys are url-encoded.
* json - json format output with values based64-encoded. keys are string.
* delete - delete matching key(s)
* List specific options.
  * showvalue - show values with keys.
  * showversion - show version keys.
  * html - add html anchor to keys for web browser operation.
  * filter:REGEXP - match keys with given regular expression.
  * max:N - limit maximum number of keys to list.
  * maxscan:N - limit maximum number of keys to scan.

Default Configuration
=====================
```
[global]
PidFile = /var/tmp/omdbd.pid
ErrorLog = /usr/local/omdb/logs/error.log
Versioning = true
MaxKeySize = 255
MaxValueSize = 0
StatusUrl = /status
QueryOption = binary,max:1000,maxscan:100000

[database]
Directory = /usr/local/omdb/db
AutoCreate = true # whether to create database automatically
CacheSize = 10000

[httplistener]
Addr = 0.0.0.0:8081     # listen address in "address:port" format
Protocol = http         # http or https
CertFile =              # certificate file for https protocol
KeyFile =               # private key file for https protocol
AccessLog = /usr/local/omdb/logs/access-20060102150405.log
LogRotate = 3600
```
