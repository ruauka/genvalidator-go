package internal

// Config - пути файлов и папок.
type Config struct {
	requestPath  string
	validatePath string
	testingPath  string
	errPath      string
}

// NewConfig - пути файлов и папок.
func NewConfig(mode string) Config {
	if mode == "debug" {
		return Config{
			requestPath:  "example/request/request.go",
			validatePath: "example/request/validate.go",
			testingPath:  "example/request/validate_test.go",
			errPath:      "example/errors",
		}
	}

	return Config{
		requestPath:  "request.go",
		validatePath: "validate.go",
		testingPath:  "validate_test.go",
		errPath:      "../errors/errors.go",
	}
}
