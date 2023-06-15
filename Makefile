VERSION=$(shell git describe --tags --always --dirty)

run:
	VERSION=${VERSION} nodemon -e go,html --signal SIGTERM --exec 'go run . || exit 1'
tw:
	tailwindcss -i input.css -o static/output.css --watch
