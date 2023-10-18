package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"text/template"

	"genvalidator/internal/templates"
)

func Execute() {
	// зачитывание и парсинг структуры в строку
	s := readStruct()

	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "", s, parser.ParseComments)
	if err != nil {
		log.Fatalf("ast file parse error: %s", err)
	}

	// проверка на необходимость перезаписи файла. Перезаписывается шапка.
	if needRewriteFile() {
		rewriteFile()
	}

	file := fileOpen()

	ast.Inspect(astFile, func(node ast.Node) bool {
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// структуры
		structs, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// поля
		for _, field := range structs.Fields.List {
			if field.Tag == nil {
				continue
			}

			for _, tag := range strings.Split(field.Tag.Value, " ") {
				// фильтр по тегу "validate"
				if !strings.HasPrefix(tag, "validate") {
					continue
				}
				// поля для шаблона
				d := templates.New(typeSpec.Name.String(), field.Names[0].String())
				// создание шаблона
				t := template.Must(template.New("require").Parse(templates.Require()))
				// дозапись шаблона в файл validate.go
				err = t.Execute(file, d)
				if err != nil {
					log.Fatalln("execute: ", err)
				}
			}
		}

		return true
	})

	file.Close()

	//ast.Print(fset, astFile)
}
