#!/bin/bash
#
# Note that this script installs the packages in this project directory,
# not in your $GOPATH. This is intended to keep the project cleanly outside
# of your personal golang dev environment.

PACKAGES="
	gopkg.in/gcfg.v1
	github.com/jmhodges/levigo
"

print_usage_and_exit() {
	cat << __EOF__
Usage:
    $0 {setup|install|clean}
__EOF__
        exit $1
}


do_setup() {
	cd src
	[[ ! -d leveldb ]] && git clone https://github.com/google/leveldb.git && (cd leveldb; git checkout v1.18; make)
	cd ../

	. env.sh
	for PACKAGE in $PACKAGES; do
		[[ -d src/$PACKAGE ]] && continue
		echo -n "Getting $PACKAGE... "
		CGO_CFLAGS="-I../../../leveldb/include/" go get -u $PACKAGE
		echo ""
	done
}

do_install() {
	[[ -f /usr/lib64/libleveldb.so ]] || [[ -f /usr/lib/libleveldb.so ]] && echo "leveldb already installed." && return
	[[ $UID != 0 ]] && echo "must be root permission to install" && return
	for DIR in /usr/lib64 /usr/lib; do
		[[ -d $DIR ]] && echo "Installing leveldb into $DIR/" && cp -v src/leveldb/libleveldb.* $DIR/ && break
	done
}

do_uninstall() {
	[[ $UID != 0 ]] && echo "must be root permission to remove" && return
	for DIR in /usr/lib64 /usr/lib; do
		[[ -f $DIR/libleveldb.so ]] && echo "Removing leveldb from $DIR/" && rm -v $DIR/libleveldb.*
	done
}

do_clean() {
	CLEAN_DIRS="pkg src/github.com src/gopkg.in src/leveldb"
	for DIR in $CLEAN_DIRS; do
		[[ -d $DIR ]] && echo "Removing $DIR... " && rm -rf $DIR
	done
}

[[ $# < 1 ]] && print_usage_and_exit 1
while [[ $# > 0 ]]; do
        option="$1"
        case $option in
        --help|-h)
                print_usage_and_exit 0;;
        setup)
                do_setup;;
        install)
                do_install;;
        uninstall)
                do_uninstall;;
        clean)
                do_clean;;
        *)
                error_exit "Unknown option: $option";;
        esac
        shift
done

exit 0

