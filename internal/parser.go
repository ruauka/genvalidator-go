package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/fatih/structtag"

	"genvalidator/internal/templates"
)

var (
	greaterThen = []string{"GreaterThen", "больше", "greater"}
	lessThen    = []string{"LessThen", "меньше", "less"}
)

func Execute() {
	// зачитывание и парсинг структуры в строку
	s := readStruct()

	fileSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fileSet, "", s, parser.ParseComments)
	if err != nil {
		log.Fatalf("ast file parse error: %s", err)
	}

	// проверка на необходимость перезаписи файла validate.go. Перезаписывается шапка.
	if needRewriteFile("validation/validate.go") {
		rewriteFile("validation/validate.go", templates.HeadValidate)
	}

	// проверка на необходимость перезаписи файла validate.go. Перезаписывается шапка.
	if needRewriteFile("validation/errors.go") {
		rewriteFile("validation/errors.go", templates.HeadErrors)
	}

	fileValidate := fileOpen("validation/validate.go")
	fileErr := fileOpen("validation/errors.go")

	errTemp := templates.NewTemplateError()

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

			// проверка на слайс (флаг слайса)
			var isArr bool
			// проверка на слайс
			_, ok := field.Type.(*ast.ArrayType)
			if ok {
				isArr = true
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
					temp          = templates.NewFunctionsTemplate()
					isFirstConcat = true
					index         int
					counter       = 1
				)

				for _, tags := range tagsParsed {
					if isFirstConcat {
						index = 172
					} else {
						index = len(temp.Buffer) - 17
					}
					switch {
					case len(tags) == 1 && tags[0] == "rq":
						// проверка на наличие уже созданной ошибки и добавление шаблона ошибки в буфер
						templates.AddErrTemplateToBuffer("rq", templates.RequireErr(), &counter, errTemp)

						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						// создание шаблона validate
						temp.Concat(templates.Require(), index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "lt":
						var (
							// создание шаблона ошибки
							errTemplate = templates.CreateErrorTemplate(greaterThen, tags[1])
							// получение названия переменной ошибки
							errVarName = strings.Split(errTemplate, " ")[5]
						)
						// проверка на наличие уже созданной ошибки и добавление шаблона ошибки в буфер
						templates.AddErrTemplateToBuffer(tags[0]+tags[1], errTemplate, &counter, errTemp)

						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.LessThan = tags[1]
						// создание шаблона
						var te string
						if isArr {
							te = templates.LessThanSl(errVarName)
						} else if isRqTag(tagsParsed) {
							te = templates.LessThanPtr(errVarName)
						} else {
							te = templates.LessThan(errVarName)
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "gt":
						var (
							// создание шаблона ошибки
							errTemplate = templates.CreateErrorTemplate(lessThen, tags[1])
							// получение названия переменной ошибки
							errVarName = strings.Split(errTemplate, " ")[5]
						)

						// проверка на наличие уже созданной ошибки и добавление шаблона ошибки в буфер
						templates.AddErrTemplateToBuffer(tags[0]+tags[1], errTemplate, &counter, errTemp)

						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.GreaterThan = tags[1]
						// создание шаблона
						var te string
						if isArr {
							te = templates.GreaterThanSl(errVarName)
						} else if isRqTag(tagsParsed) {
							te = templates.GreaterThanPtr(errVarName)
						} else {
							te = templates.GreaterThan(errVarName)
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					}
				}

				// создание и запись шаблона в validate.go
				templateExecute("validate.go", temp.Buffer, fileValidate, templFields)
				// очистка буфера шаблона FunctionsTemplate
				temp.Reset()
			}
		}

		return true
	})

	// создание и запись шаблона в errors.go
	templateExecute("errors.go", errTemp.Buffer, fileErr, templates.TemplateFields{})

	fileValidate.Close()
	fileErr.Close()

	//ast.Print(fset, astFile)
}

// templateExecute - создание и запись шаблона в файл.go.
func templateExecute(name, buffer string, file *os.File, fields templates.TemplateFields) {
	// создание шаблона
	p, err := template.New(name).Parse(buffer)
	if err != nil {
		log.Fatalln("TemplateExecute(). Ошибка создания шаблона: ", err)
	}
	// запись шаблона в файл
	if err = template.Must(p, err).Execute(file, fields); err != nil {
		log.Fatalln("TemplateExecute(). Ошибка записи шаблона в файл.go: ", err)
	}
}
