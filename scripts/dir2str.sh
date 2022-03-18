#!/bin/bash

delim="#"

directory=$(echo ${PWD##*/})
outname="${directory}.tar.bz2"
tar cjf "$outname" *
base64Str=$(base64 "$outname" | tr -d '\n')
rm "$outname"
echo "$(echo -ne "$directory" | base64)${delim}$base64Str"
