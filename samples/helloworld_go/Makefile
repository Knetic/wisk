default: build
all: package

export GOPATH=$(CURDIR)/
export GOBIN=$(CURDIR)/.temp/

init: clean
	go get ./...

build: init
	go build -o ./.output/${{=project.name=}} .

test:
	go test
	go test -bench=.

clean:
	rm -rf ./.output/

dist: build test
	export GOOS=linux
	export GOARCH=386
	go build -o ./.output/${{=project.name=}}32 .

	export GOOS=darwin
	export GOARCH=amd64
	go build -o ./.output/${{=project.name=}}_osx .

	export GOOS=windows
	export GOARCH=amd64
	go build -o ./.output/${{=project.name=}}.exe .



package: dist

ifeq ($(shell which fpm), )
	@echo "FPM is not installed, no packages will be made."
	@echo "https://github.com/jordansissel/fpm"
	exit 1
endif

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v 1.0 \
		-n ${{=project.name=}} \
		./.output/${{=project.name=}}=/usr/local/bin/${{=project.name=}}

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v 1.0 \
		-n ${{=project.name=}} \
		-a i686 \
		./.output/${{=project.name=}}32=/usr/local/bin/${{=project.name=}}

	mv ./*.deb ./.output/

	# rpm
	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v 1.0 \
		-n ${{=project.name=}} \
		./.output/${{=project.name=}}=/usr/local/bin/${{=project.name=}}
	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v 1.0 \
		-n ${{=project.name=}} \
		-a i686 \
		./.output/${{=project.name=}}32=/usr/local/bin/${{=project.name=}}

	mv ./*.rpm ./.output/
