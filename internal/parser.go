package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
)

var checkMap = make(map[string]struct{})

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
		rewriteFile("validation/validate.go", HeadValidate)
	}

	// проверка на необходимость перезаписи файла validate.go. Перезаписывается шапка.
	if needRewriteFile("validation/errors.go") {
		rewriteFile("validation/errors.go", HeadErrors)
	}

	fileValidate := fileOpen("validation/validate.go")
	fileErr := fileOpen("validation/errors.go")

	errTemp := NewTemplateError()

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
					templFields   TemplateFields
					tValid        *template.Template
					temp          = NewTemplate()
					isFirstConcat = true
					index         int
					counter       = 1
				)

				for _, tags := range tagsParsed {
					if isFirstConcat {
						index = 65
					} else {
						index = len(temp.Buffer) - 17
					}
					switch {
					case len(tags) == 1 && tags[0] == "rq":
						// проверка на наличие уже созданной ошибки
						if _, ok := checkMap["rq"]; !ok {
							// создание шаблона error
							errTemp.Concat(RequireErr(), 6, counter)
							counter++
							checkMap["rq"] = struct{}{}
						}
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						// создание шаблона validate
						temp.Concat(Require(), index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "lt":
						var errStr = createErr([]string{"GreaterThen", "больше", "greater"}, tags[1])
						// проверка на наличие уже созданной ошибки
						_, ok := checkMap[tags[0]+tags[1]]
						if !ok {
							// создание шаблона error
							errTemp.Concat(errStr, 6, counter)
							counter++
							checkMap[tags[0]+tags[1]] = struct{}{}
						}
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.LessThan = tags[1]
						// создание шаблона
						var te string
						if isArr {
							te = LessThanSl(strings.Split(errStr, " ")[5])
						} else if isPtr(tagsParsed) {
							te = LessThanPtr(strings.Split(errStr, " ")[5])
						} else {
							te = LessThan(strings.Split(errStr, " ")[5])
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					case tags[0] == "gt":
						var errStr = createErr([]string{"LessThen", "меньше", "less"}, tags[1])
						// проверка на наличие уже созданной ошибки
						_, ok := checkMap[tags[0]+tags[1]]
						if !ok {
							// создание шаблона error
							errTemp.Concat(errStr, 6, counter)
							counter++
							checkMap[tags[0]+tags[1]] = struct{}{}
						}
						// поля для шаблона
						templFields.StructName = typeSpec.Name.String()
						templFields.StrategyFieldName = field.Names[0].String()
						templFields.JsonFieldName = jsonName.Name
						templFields.GreaterThan = tags[1]
						// создание шаблона
						var te string
						if isArr {
							te = GreaterThanSl(strings.Split(errStr, " ")[5])
						} else if isPtr(tagsParsed) {
							te = GreaterThanPtr(strings.Split(errStr, " ")[5])
						} else {
							te = GreaterThan(strings.Split(errStr, " ")[5])
						}
						temp.Concat(te, index, isFirstConcat)
						isFirstConcat = false
					}
				}

				var err error
				// создание шаблона validate.go
				p, err := template.New("").Parse(temp.Buffer)
				if err != nil {
					log.Fatalln(111, err)
				}

				tValid = template.Must(p, err)

				err = tValid.Execute(fileValidate, templFields)
				if err != nil {
					log.Fatalln("execute: ", err)
				}

				temp.Reset()

			}
		}

		return true
	})

	// создание шаблона error.go
	e, err := template.New("").Parse(errTemp.Buffer)
	if err != nil {
		log.Fatalln(111, err)
	}

	var tErr *template.Template
	tErr = template.Must(e, err)

	err = tErr.Execute(fileErr, "")
	if err != nil {
		log.Fatalln("execute: ", err)
	}

	fileValidate.Close()
	fileErr.Close()

	//ast.Print(fset, astFile)
}
