package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/fatih/structtag"

	"genvalidator/internal/rules"
	"genvalidator/internal/templates"
)

// Paths - пути файлов и папок.
type Paths struct {
	requestPath  string
	validatePath string
	testingPath  string
	errPath      string
}

// New - пути файлов и папок.
func New(mode string) Paths {
	if mode == "debug" {
		return Paths{
			requestPath:  "validation/request/request.go",
			validatePath: "validation/request/validate.go",
			testingPath:  "validation/request/validate_test.go",
			errPath:      "validation/errors",
		}
	}

	return Paths{
		requestPath:  "request.go",
		validatePath: "validate.go",
		testingPath:  "validate_test.go",
		errPath:      "../errors/errors.go",
	}
}

const (
	// название правил и тега
	require     = "rq"
	lessThan    = "lt"
	greaterThan = "gt"
)

var (
	// создание буффера для шаблона "ошибка"
	errBuffer = templates.NewTemplateError()
	// создание буффера для шаблона функции "Validate"
	invokeBuffer = templates.NewInvokeTemplate()
	// создание буффера для шаблона "тестирование"
	testingBuffer = templates.NewTestingTemplate()

	// мапа для проверки уже имеющейся "ошибки"
	isErrExists = make(map[string]struct{})
	// мапа для проверки уже имеющегося вызова функции правил
	isFuncCall = make(map[string]struct{})
	// мапа для проверки уже имеющейся тестовой функции
	isTestingExists = make(map[string]struct{})

	// флаг для конкатенации буффера "тестирование"
	isFirstConcatTesting = true
)

func Execute(mode string) {
	paths := New(mode)
	// зачитывание и парсинг файла со структурами в строку
	s := readStruct(paths.requestPath)

	fs := token.NewFileSet()
	// дерево ast
	astFile, err := parser.ParseFile(fs, "", s, parser.ParseComments)
	if err != nil {
		log.Fatalf("ast file parse error: %s", err)
	}

	// определение путей для папки errors и файла errors.go для дебага и пром режима
	var errPathDir, errPathFile string

	if mode == "debug" {
		errPathDir, errPathFile = paths.errPath, path.Join(paths.errPath, "errors.go")
	} else {
		errPathDir, errPathFile = paths.errPath[:9], paths.errPath
	}

	// проверка наличия папки errors, если нет, то создать
	if err := os.Mkdir(errPathDir, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Создание файла.go с нужной шапкой validate.go | errors.go | validate_test.go
	createFileWithTemplate(paths.validatePath, templates.HeadValidate)
	createFileWithTemplate(errPathFile, templates.HeadErrors)
	createFileWithTemplate(paths.testingPath, func() string { return "" })

	// открытие файлов на дозапись validate.go | errors.go | validate_test.go
	var (
		fileValidate = fileOpenAppendMode(paths.validatePath)
		fileErr      = fileOpenAppendMode(errPathFile)
		fileTesting  = fileOpenAppendMode(paths.testingPath)
	)

	// проход по нодам дерева ast
	ast.Inspect(astFile, func(node ast.Node) bool {
		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// получение структур из ast
		structs, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		// получение полей
		for _, field := range structs.Fields.List {
			// проверка отсутствия тега
			if field.Tag == nil {
				continue
			}
			// флаг типа поля - слайс
			var isArr bool
			// проверка на типа поля - слайс
			if _, ok := field.Type.(*ast.ArrayType); ok {
				isArr = true
			}
			// парсингвсех тегов
			tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
			if err != nil {
				log.Fatalf("structtag.Parse. Ошибка парсинга тегов: %s", err)
			}
			// получение названия поля из json
			jsonFieldName, err := tags.Get("json")
			if err != nil {
				log.Fatalf("tags.Get(\"json\"). Ошибка получения названия поля из json: %s", err)
			}
			// получение тэгов
			for _, tag := range tags.Tags() {
				// фильтр по тегу "validate"
				if tag.Key != "validate" {
					continue
				}

				var (
					// парсинг правил тега "validate" в [][]string | [[rq]] | [[rq] [lt 10]]
					rulesParsed = validateTagParse(tag.Value())
					// поля для шаблона "функция"
					funcFields templates.TemplateFields
					// создание буффера для шаблона "функция"
					funcBuffer = templates.NewFunctionsTemplate()
					// флаг для конкатенации буффера "функция" | "ошибка" | функции "Validate"
					isFirstConcatFunc, isFirstConcatErr, isFirstConcatValidation = true, true, true
				)

				// палучение правил тега
				for _, tagRules := range rulesParsed {
					// расчет отступа для взятия подстроки для конкатенации буффера
					indent := indentConcat(isFirstConcatFunc, funcBuffer.Buffer)
					// основная логика
					switch {
					// обработка правила "rq"
					case len(tagRules) == 1 && tagRules[0] == require:
						rules.Require(
							&isFirstConcatFunc, &isFirstConcatErr, &isFirstConcatValidation, &isFirstConcatTesting,
							isErrExists, isTestingExists, isFuncCall,
							errBuffer, funcBuffer,
							invokeBuffer, testingBuffer,
							&funcFields,
							typeSpec, field, jsonFieldName,
							indent,
						)
					// обработка правила "lt"
					case tagRules[0] == lessThan:
						rules.LessThan(
							&isFirstConcatFunc, &isFirstConcatErr, &isFirstConcatValidation, &isFirstConcatTesting,
							isErrExists, isTestingExists, isFuncCall,
							errBuffer, funcBuffer,
							invokeBuffer, testingBuffer,
							&funcFields,
							typeSpec, field, jsonFieldName,
							indent,
							isArr,
							tagRules,
							rulesParsed,
						)
					// обработка правила "gt"
					case tagRules[0] == greaterThan:
						rules.GreaterThan(
							&isFirstConcatFunc, &isFirstConcatErr, &isFirstConcatValidation, &isFirstConcatTesting,
							isErrExists, isTestingExists, isFuncCall,
							errBuffer, funcBuffer,
							invokeBuffer, testingBuffer,
							&funcFields,
							typeSpec, field, jsonFieldName,
							indent,
							isArr,
							tagRules,
							rulesParsed,
						)
					}
				}
				// создание и запись шаблона "функции" в validate.go
				templateExecute("validate.go", funcBuffer.Buffer, fileValidate, funcFields)
				// очистка буфера шаблона "функции" FunctionTemplate
				funcBuffer.Reset()
			}
		}

		// переход следующую ноду дерева ast
		return true
	})

	// создание и запись шаблона функции "Validate" в validate.go
	templateExecute("validate.go", invokeBuffer.Buffer, fileValidate, templates.TemplateFields{})
	// создание и запись шаблона "ошибки" в errors.go
	templateExecute("errors.go", errBuffer.Buffer, fileErr, templates.TemplateFields{})
	// создание и запись шаблона "тестирование" в validate_test.go
	templateExecute("errors.go", testingBuffer.Buffer, fileTesting, templates.TemplateFields{})

	defer func() { _ = fileValidate.Close() }()
	defer func() { _ = fileErr.Close() }()
	defer func() { _ = fileTesting.Close() }()

	// ast.Print(fs, astFile)
}

// indentConcat - расчет отступа для взятия подстроки для конкатинации буффера.
func indentConcat(isFirstConcat bool, buffer string) int {
	if isFirstConcat {
		return 172
	}

	return len(buffer) - 17
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

// validateTagParse - парсинг правил тега "validate".
func validateTagParse(rules string) [][]string {
	var (
		splited = strings.Split(rules, ",")
		trimmed = make([][]string, len(splited))
	)

	for idx, sl := range splited {
		splited := strings.Split(sl, "=")
		for _, val := range splited {
			trimmed[idx] = append(trimmed[idx], strings.TrimSpace(val))
		}
	}

	// постановка тега "rq" на первое место, на случай если последовательность правил начинается не с "rq"
	if rqCheck(trimmed) {
		var (
			trimmedSorted = make([][]string, 0, len(trimmed))
			index         int
		)

		for idx, sl := range trimmed {
			if len(sl) == 1 && sl[0] == require {
				trimmedSorted = append(trimmedSorted, sl)
				index = idx
				break
			}
		}

		trimmedSorted = append(trimmedSorted, trimmed[:index]...)
		trimmedSorted = append(trimmedSorted, trimmed[index+1:]...)

		return trimmedSorted
	}

	return trimmed // [[rq]] или [[lt 10]]
}

// rqCheck - проверка наличия правила "rq".
func rqCheck(rules [][]string) bool {
	for _, sl := range rules {
		if len(sl) == 1 && sl[0] == require {
			return true
		}
	}

	return false
}
