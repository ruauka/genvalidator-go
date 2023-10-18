package internal

import (
	"io"
	"log"
	"os"
	"strings"

	"genvalidator/internal/templates"
)

func readStruct() string {
	f, err := os.Open("example.go")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	b := new(strings.Builder)
	io.Copy(b, f)

	return b.String()
}

func fileOpen() *os.File {
	// открытие файла в режиме дозаписывания
	file, err := os.OpenFile("validate.go", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("os.OpenFile err: %s", err)
	}

	return file
}

func needRewriteFile() bool {
	file, err := os.Stat("validate.go")
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}

	if file.Size() != 0 {
		return true
	}

	return false
}

// rewriteFile - перезапись файла.
func rewriteFile() {
	err := os.WriteFile("validate.go", []byte(templates.Head()), 0644)
	if err != nil {
		log.Fatalf("os.WriteFile err: %s", err)
	}

}

// validateTagParse - парсинг правил.
func validateTagParse(tag string) [][]string {
	var (
		splited = strings.Split(tag, ",")
		trimmed = make([][]string, len(splited))
	)

	for idx, sl := range splited {
		splited := strings.Split(sl, "=")
		for _, val := range splited {
			trimmed[idx] = append(trimmed[idx], strings.TrimSpace(val))
		}
	}

	return trimmed // [[rq] [lt 10]]
}

func isPtr(tags [][]string) bool {
	for _, tag := range tags {
		if tag[0] == "rq" {
			return true
		}
	}

	return false
}
