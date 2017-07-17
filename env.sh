if [ -z "$GOPATH" ]; then
	export GOPATH=$(pwd)
else
	export GOPATH=$GOPATH:$(pwd)
fi

export GOPATH=$GOPATH:/Users/shaozhenyu/work/git/utils
