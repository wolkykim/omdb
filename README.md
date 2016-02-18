IN-PROGRESS OF DEVELOPMENT
==========================

OmDB (OhMyDB)
=============

OmDB is a persistent key/value store. OmDB uses LevelDB as backend storre and provide REST APIs.

* OmDB is simple as single binary executable is all you need to install.
* OmDB doesn't require special client but standard HTTP(HTTPS) protocol.
* OmDB can be accessed using only curl and OmDB is shell program friendly.
* OmDB provides sorted key listing and key search.
* OmDB provides iteration over lists.
* OmDB provides data versioning.


How to Build
============
```
$ git clone https://github.com/wolkykim/omdb.git
$ cd omdb
$ source env.sh
$ make updatepkg
$ sudo make install-leveldb-lib
$ make
$ src/omdbd/omdbd -c etc/omdbd.conf -d
```

APIs
====


## LIST API

DEMO: http://omdb.theaddr.com:8081/test/

List all keys in sorted order that are prefixed by a given key or within the given range.
```
GET /db/key/ HTTP/1.1

200 OK
Content-Type: text/plain
Content-Length: 100
x-key: key
x-size: 3
x-format: text
x-encoding: url
x-delimiter: (Not shown if not specified in the request)
x-keyfilter: (Not shown if not specified in the request)
x-rangekey: (Now shown if not specified in the request)
x-iterator: (Not shown if no iteration is required)

keyA=value
keyB=value
keyC=value
```

For the next iteration. (Also can be given as an option)
```
GET /db/key/#iterator HTTP/1.1
```

Options can be given in 2 ways
```
GET /db/key/?option1=value&option2=value&... HTTP/1.1

Or

GET /db/key/ HTTP/1.1
x-options: option1=value,option2=value,...
```

List of Options
* maxkey : Maximum number of keys (default: 1000)
* format : text | html | xml | json (default: text)
* keyonly : true | false (default false)
* encoding : url | hex | base64 (default: url)
* rangekey : last key for range request (if given perform until the specified key range, if not given perform prefix listing)
* delimiter : character for the boundary (ex: "/")
* keyfilter : regular expression to filter keys.
* iterator : Used for iteration

Response Codes
* 200 : Ok.
* 400 : Bad request.
* 500 : Internal server error

Response Headers
* Content-Type
* Content-Length
* x-key
* x-size
* x-format
* x-encoding
* x-iterator: Will be set if list is truncated by maxkey.


## GET API

DEMO: http://omdb.theaddr.com:8081/test/hello

Get the value of a key.
```
GET /db/key HTTP/1.1

200 OK
Content-Type: application/octet-stream
Content-Length: 32
x-key: /db/key
x-encoding: raw

(raw data)
```

List of Options
* encoding : raw | url | hex | base64 | gzip (default: raw)
* contenttype : override default content-type to user defined (default: application/octet-stream for 'raw' and 'gzip' incoding, plain/text for other 'encoding')

Response Codes
* 200 : Ok.
* 400 : Bad request.
* 404 : Key not found.
* 500 : Internal server error

Response Headers
* Content-Type
* Content-Length
* x-key
* x-options


## PUT API

DEMO: http://omdb.theaddr.com:8081/test/testkey?v=test_value

Put a key
```
PUT /db/key HTTP/1.1
Content-Length: 32

(raw data)

204 OK
Content-Length: 32
x-key: /db/key
```

List of Options
* encoding : raw | gzip (default: raw)

Response Codes
* 204 : Successfully stored.
* 400 : Bad request.
* 404 : Key not found.
* 500 : Internal server error

Response Headers
* Content-Length
* x-key

## DELETE API

Delete a key
```
DELETE /db/key HTTP/1.1

204 OK
Content-Length: 32
x-key: /db/key
```

Delete all keys starting with key prefix.
```
DELETE /db/key_prefix/ HTTP/1.1

200 OK
Content-Length: 32
x-key: /db/key_prefix
x-size: 3
x-encoding: url
x-iterator:

200 /db/key_prefix_A
200 /db/key_prefix_B
500 /db/key_prefix_C
```

Delete multiple keys.
```
DELETE / HTTP/1.1
x-key: /db/key_A
x-key: /db/key_B, /db/key_C

200 OK
Content-Length: 32
x-key: /
x-size: 3
x-encoding: url
x-iterator:

200 /db/key_A
404 /db/key_B
200 /db/key_C
```


List of Options
* maxkey : Maximum number of keys to delete (default: 1000)
* format : text | html | xml | json (default: text)
* encoding url | hex | base64 (default: url)
* iterator : Used for iteration

Response Codes
* 200 or 204 : Successfully deleted.
* 400 : Bad request.
* 404 : Key not found.
* 500 : Internal server error

Response Headers
* Content-Length
* x-key

Configuration
=============
```
[global]
PidFile = /var/tmp/omdbd.pid
ErrorLog = /usr/local/omdbd/logs/error.log

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
#StatusUrl = /status

[limit]
MaxKeySize = 255
MaxValueSize = 0
ListSize = 1000
```
