package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"text/template"

	"github.com/fatih/structtag"

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
			// парсингвсех тегов
			tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
			if err != nil {
				log.Fatalf("structtag.Parse. Tags parse error: %s", err)
			}

			jsonName, err := tags.Get("json")
			if err != nil {
				panic(err)
			}

			for _, tag := range tags.Tags() {

				// фильтр по тегу "validate"
				if tag.Key != "validate" {
					continue
				}

				tagsParsed := validateTagParse(tag.Value())
				// [[rq]]
				// [[rq] [lt 10]]

				var (
					templFields   templates.TemplateFields
					t             *template.Template
					temp          = templates.NewTemplate()
					isFirstConcat = true
					index         int
				)

				for _, tags := range tagsParsed {
					if isFirstConcat {
						index = 65
					} else {
						index = len(temp.Buffer) - 17
					}
					switch {
					case len(tags) == 1 && tags[0] == "rq":
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						// создание шаблона
						temp.Concat(templates.Require(), index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "lt":
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.LessThan = tags[1]
						// создание шаблона
						var te string
						if isPtr(tagsParsed) {
							te = templates.LessThanPtr()
						} else {
							te = templates.LessThan()
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "gte":
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.GreaterThanOrEq = tags[1]
						// создание шаблона
						var te string
						if isPtr(tagsParsed) {
							te = templates.GreaterThanOrEqPtr()
						} else {
							te = templates.GreaterThanOrEq()
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					}
				}

				var err error
				// создание шаблона
				p, err := template.New("").Parse(temp.Buffer)
				if err != nil {
					log.Fatalln(111, err)
				}

				t = template.Must(p, err)

				err = t.Execute(file, templFields)
				if err != nil {
					log.Fatalln("execute: ", err)
				}

				temp.Reset()
			}
		}

		return true
	})

	file.Close()

	//ast.Print(fset, astFile)
}
