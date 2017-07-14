#!/bin/bash

function getdir(){
	for element in `ls $1`; do
		file=$1"/"$element
		if [ -d $file ]; then
			getdir $file
		else
			echo "$file" | grep "glide.yaml" > /dev/null
			if [ $? -eq 0 ]; then
				prefix=${file%/*}
				cd "$prefix"
				glide update
			fi
		fi
	done
}

dir="./src"
target="glide.yaml"
getdir $dir

