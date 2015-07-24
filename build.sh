export GOPATH="$(pwd)"
export GOOS=linux
export GOARCH=amd64

function packageDeb() {

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

	mv ./*.deb ./.output/
}

function buildExecutable() {

  go get ./...
  go build -o ./.output/wisk .
}

buildExecutable

if [[ "$1" != "package" ]];
then
	exit 0
fi

packageDeb
