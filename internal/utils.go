package internal

import (
	"io"
	"log"
	"os"
	"strings"
)

func readStruct() string {
	f, err := os.Open("validation/request.go")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	b := new(strings.Builder)
	io.Copy(b, f)

	return b.String()
}

func fileOpen(path string) *os.File {
	// открытие файла в режиме дозаписывания
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("os.OpenFile err: %s", err)
	}

	return file
}

func needRewriteFile(path string) bool {
	file, err := os.Stat(path)
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
func rewriteFile(path string, fn func() string) {
	err := os.WriteFile(path, []byte(fn()), 0644)
	if err != nil {
		log.Fatalf("os.WriteFile err: %s", err)
	}

}
