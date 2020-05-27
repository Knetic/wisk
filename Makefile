default: build
all: package

export GOPATH=$(CURDIR)/
export GOBIN=$(CURDIR)/.temp/

init: clean
	go get ./...

build: init
	go build -o ./.output/wisk .

test:
	go test
	go test -bench=.

integrate: build
	@./.output/wisk -a ./samples/helloworld_ruby/
	./.output/wisk helloworld_ruby ./.populatedSample

clean:
	@rm -rf ./.output/

dist: build test

	export GOOS=linux; \
	export GOARCH=amd64; \
	go build -o ./.output/wisk64 .

	export GOOS=linux; \
	export GOARCH=386; \
	go build -o ./.output/wisk32 .

	export GOOS=darwin; \
	export GOARCH=amd64; \
	go build -o ./.output/wisk_osx .

	export GOOS=windows; \
	export GOARCH=amd64; \
	go build -o ./.output/wisk.exe .

package: dist

ifeq ($(shell which fpm), )
	@echo "FPM is not installed, no packages will be made."
	@echo "https://github.com/jordansissel/fpm"
	@exit 1
endif

ifeq ($(WISK_VERSION), )

	@echo "No 'WISK_VERSION' was specified."
	@echo "Export a 'WISK_VERSION' environment variable to perform a package"
	@exit 1
endif

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(WISK_VERSION) \
		-n wisk \
		./.output/wisk64=/usr/local/bin/wisk \
		./docs/wisk.7=/usr/share/man/man7/wisk.7 \
		./autocomplete/wisk=/etc/bash_completion.d/wisk

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(WISK_VERSION) \
		-n wisk \
		-a i686 \
		./.output/wisk32=/usr/local/bin/wisk \
		./docs/wisk.7=/usr/share/man/man7/wisk.7 \
		./autocomplete/wisk=/etc/bash_completion.d/wisk

	@mv ./*.deb ./.output/

	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(WISK_VERSION) \
		-n wisk \
		./.output/wisk64=/usr/local/bin/wisk \
		./docs/wisk.7=/usr/share/man/man7/wisk.7 \
		./autocomplete/wisk=/etc/bash_completion.d/wisk

	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(WISK_VERSION) \
		-n wisk \
		-a i686 \
		./.output/wisk32=/usr/local/bin/wisk \
		./docs/wisk.7=/usr/share/man/man7/wisk.7 \
		./autocomplete/wisk=/etc/bash_completion.d/wisk

	@mv ./*.rpm ./.output/

dockerTest:
ifeq ($(shell which docker), )
	@echo "Docker is not installed."
	@exit 1
endif

containerized_build: dockerTest

	docker run \
		--rm \
		-v "$(CURDIR)":"/srv/build":rw \
		-e ALT_VERSION=$(SEEDSTATS_VERSION) \
		golang:1.13 \
		bash -c \
		"cd /srv/build; make build"