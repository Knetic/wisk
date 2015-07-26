#!/bin/bash

export GOPATH="$(pwd)"

function package() {

	if [ "$(which fpm)" == "" ];
	then
		echo "FPM is not installed, no packages will be made."
		echo "https://github.com/jordansissel/fpm"
		exit 1
	fi

	rm -f ./.output/*.deb

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v 1.0 \
		-n wisk \
		./.output/wisk=/usr/local/bin/wisk

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v 1.0 \
		-n wisk \
		-a i686 \
		./.output/wisk32=/usr/local/bin/wisk

	mv ./*.deb ./.output/

	# rpm
	rm -f ./.output/*.rpm

	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v 1.0 \
		-n wisk \
		./.output/wisk=/usr/local/bin/wisk
	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v 1.0 \
		-n wisk \
		-a i686 \
		./.output/wisk32=/usr/local/bin/wisk

	mv ./*.rpm ./.output/
}

go get ./...

export GOOS=linux
export GOARCH=amd64
go build -o ./.output/wisk

if [[ "$1" != "package" ]];
then
	exit 0
fi

export GOOS=linux
export GOARCH=386
go build -o ./.output/wisk32 .

export GOOS=darwin
export GOARCH=amd64
go build -o ./.output/wisk_osx .

export GOOS=windows
export GOARCH=amd64
go build -o ./.output/wisk.exe .

package
