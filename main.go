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
	defer func() {
		if err := recover(); err != nil {
			log.Println("empty request or errors path", err)
		}
	}()

	args := os.Args

	request := args[1] // validation/request
	errors := args[2]  // validation/errors

	internal.Execute(
		request+requestFilePath,
		request+validationFilePath,
		request+testingFilePath,
		errors+errorsFilePath,
	)
}
