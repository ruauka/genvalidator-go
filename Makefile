build-amd64:
	GOOS=darwin GOARCH=amd64 go build -o genvalidate main.go

build-arm64:
	GOOS=darwin GOARCH=arm64 go build -o genvalidate main.go #m2

start:
	go run main.go validation/request validation/errors

startb:
	./genvalidate validation/request validation/errors