#!/bin/bash

delim="#"

if [ $# -ne 1 ]; then
	echo "Usage: $0 str"
	exit
fi

tmpFile="tmp.san.str2dir.tar.bz2"
str="$1"
dirname="$(echo $str | awk -F"$delim" '{print $1}' | base64 -d)"
hashStr=$(echo $str | awk -F"$delim" '{print $2}')

mkdir "$dirname"
cd "$dirname"
echo -ne $hashStr | base64 -d >$tmpFile
tar xf $tmpFile
rm -rf $tmpFile
echo "decompressed files in directory '$dirname'"
