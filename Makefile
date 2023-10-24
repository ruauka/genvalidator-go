build:
	go build -o genvalidate main.go

start:
	go run main.go validation/request validation/errors

startb:
	./genvalidate validation/request validation/errors