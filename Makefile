run:
	nodemon -e go,html --signal SIGTERM --exec 'go run main.go || exit 1'
