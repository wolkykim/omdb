#!/bin/bash
#
# This script installs/updates the packages used in mdpd. This is called by
# build process automatically.
#
# Note that this script installs the packages in this project directory,
# not in your $GOPATH. This is intended to keep the project cleanly outside
# of your personal golang dev environment.
#
. env.sh

PACKAGES="
	gopkg.in/gcfg.v1
"

export GOPATH=$PWD
for PACKAGE in $PACKAGES; do
	echo -n "Getting $PACKAGE... "
	go get -u $PACKAGE
	echo ""
done

cd src
[[ ! -d leveldb ]] && git clone https://github.com/google/leveldb.git && (cd leveldb; git checkout v1.18; make)
cd ../
echo -n "Getting levigo... "
CGO_CFLAGS="-I../../../leveldb/include/" go get -u github.com/jmhodges/levigo
echo ""

