package internal

import (
	"io"
	"log"
	"os"
	"strings"
)

// readStruct - зачитывание структуры из файла в строчку.
func readStruct(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("readStruct() err: %s", err)
	}

	defer func() { _ = f.Close() }()

	b := new(strings.Builder)
	io.Copy(b, f)

	return b.String()
}

// fileOpenAppendMode - открытие файла в режиме на дозапись.
func fileOpenAppendMode(path string) *os.File {
	// открытие файла в режиме дозаписывания
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("os.OpenFile err: %s", err)
	}

	return file
}

// needReWriteFile - проверка на необходимость перезаписи файла.go.
func needReWriteFile(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		// проверка наличия
		if os.IsNotExist(err) {
			return true
		}
	}
	// проверяет пустой ли файл
	if file.Size() != 0 {
		return true
	}

	return false
}

// reWriteFile - перезапись файла.
func reWriteFile(path string, fn func() string) {
	if err := os.WriteFile(path, []byte(fn()), 0644); err != nil {
		log.Fatalf("rewriteFile err: %s", err)
	}
}
