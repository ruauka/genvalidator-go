package main

import (
	"log"
	"os"

	"genvalidator/internal"
)

var (
	requestFilePath    = "/request.go"
	validationFilePath = "/validate.go"
	testingFilePath    = "/validate_test.go"
	errorsFilePath     = "/errors.go"
)

func main() {
	request := os.Getenv("request") // validation/request
	if request == "" {
		log.Fatal("request is not define")
	}
	errors := os.Getenv("errors") // validation/errors
	if errors == "" {
		log.Fatal("errors is not define")
	}

	internal.Execute(
		request+requestFilePath,
		request+validationFilePath,
		request+testingFilePath,
		errors+errorsFilePath,
	)
}

// request=validation/request errors=validation/errors go run main.go
// request=pkg/request errors=pkg/errors ./validgen
