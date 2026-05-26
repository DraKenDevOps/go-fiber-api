install:
	go mod tidy

dev:
	go run main.go

dev-re:
	re go run main.go

dev-hot-win:
	npx nodemon --exec go run main.go --signal SIGINT

dev-hot-mac:
	npx nodemon --exec go run main.go --signal SIGTERM
