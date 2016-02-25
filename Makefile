################################################################################
## OmDB
##
## Copyright (c) 2016 Seungyoung Kim.
## All rights reserved.
##
## Redistribution and use in source and binary forms, with or without
## modification, are permitted provided that the following conditions are met:
##
## 1. Redistributions of source code must retain the above copyright notice,
##    this list of conditions and the following disclaimer.
## 2. Redistributions in binary form must reproduce the above copyright notice,
##    this list of conditions and the following disclaimer in the documentation
##    and/or other materials provided with the distribution.
##
## THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
## AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
## IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
## ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
## LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
## CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
## SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
## INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
## CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
## ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
## POSSIBILITY OF SUCH DAMAGE.
################################################################################

.PHONY: all build install clean package rpm
D=$(shell pwd)/rpm-build
O=--define "_topdir $(D)"
PREFIX=/usr/local/omdb

TARGETS=src/omdbd

all: build

build-deps:
	@./update_pkg.sh setup

install-deps:
	@./update_pkg.sh install

uninstall-deps:
	@./update_pkg.sh uninstall

build: build-deps
	@for TARGET in ${TARGETS}; do				\
		(export GOPATH="$${PWD}"; cd $${TARGET}; make);	\
	done

install: build install-deps
	@mkdir -p ${PREFIX}
	@mkdir -p ${PREFIX}/bin
	@mkdir -p ${PREFIX}/db
	@mkdir -p ${PREFIX}/etc
	@mkdir -p ${PREFIX}/logs
	@cp -fv src/omdbd/omdbd ${PREFIX}/bin/
	@cp -fv etc/omdbd.conf.example ${PREFIX}/etc/

clean:
	@./update_pkg.sh clean
	@for TARGET in ${TARGETS}; do				\
		(cd $${TARGET}; make clean);			\
	done

package: rpm
rpm:	clean build
	rpmbuild -bb --noclean $(O) rpm.spec
