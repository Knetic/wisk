default: build
all: package

export GOPATH=$(CURDIR)/
export GOBIN=$(CURDIR)/.temp/

FPM=$(shell which fpm)

init: clean
	go get ./...

build: init
	go build -o ./.output/wisk .

test:
	go test
	go test -bench=.

clean:
	rm -rf ./.output/

dist: build test
	export GOOS=linux
	export GOARCH=386
	go build -o ./.output/wisk32 .

	export GOOS=darwin
	export GOARCH=amd64
	go build -o ./.output/wisk_osx .

	export GOOS=windows
	export GOARCH=amd64
	go build -o ./.output/wisk.exe .


package: dist

ifndef FPM
	echo "FPM is not installed, no packages will be made."
	echo "https://github.com/jordansissel/fpm"
	exit 1
endif

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
