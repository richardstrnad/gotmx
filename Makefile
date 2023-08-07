VERSION=$(shell git describe --tags --always --dirty)

run:
	VERSION=${VERSION} nodemon -e go,html --signal SIGTERM --exec 'go run cmd/webd/webd.go || exit 1'

build:
	docker build -t richardstrnad/gotmx:${VERSION} . --build-arg VERSION=$(VERSION)

tw:
	tailwindcss -i input.css -o static/output.css --watch

test:
	VERSION="TESTVERSION" go test -v ./...
