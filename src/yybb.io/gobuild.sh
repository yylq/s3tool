#!/bin/bash
set -e

[[ -z "$GOPATH" ]] && GOPATH=$(pwd)/../../
[[ -z "$OUTPUT" ]] && OUTPUT=$GOPATH/output

if [ ! -d "$OUTPUT/bin" ]; then
   mkdir -p $OUTPUT/bin	
fi
echo $GOPATH
PACKAGE=$1
BIN=$(echo $PACKAGE|cut -d '/' -f2)
go build -o $OUTPUT/bin/$BIN $PACKAGE
