#!/usr/bin/env bash

# if the environment has been setup before clean it up
if [ $GOBIN ]; then
    PATH=$(echo $PATH | sed -e "s;\(^$GOBIN:\|:$GOBIN$\|:$GOBIN\(:\)\);\2;g")
fi

export GOPATH=`pwd`
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

