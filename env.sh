if [ -z "$GOPATH" ]; then
	export GOPATH=$(pwd)
else
	export GOPATH=$GOPATH:$(pwd)
fi

echo "glide update"
bash glide.sh
