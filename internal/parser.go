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

// название правил и теге
const (
	require     = "rq"
	lessThan    = "lt"
	greaterThan = "gt"
)

var (
	greaterThen          = []string{"GreaterThen", "больше", "greater"}
	lessThen             = []string{"LessThen", "меньше", "less"}
	isFirstConcatTesting = true // флаг для конкатенации буффера "тестирование"
	errTemplate          string
	errVarName           string
)

var (
	// мапа для проверки уже имеющейся "ошибки"
	isErrExists = make(map[string]struct{})
	// мапа для проверки уже имеющегося вызова функции правил
	isFuncCall = make(map[string]struct{})
	// мапа для проверки уже имеющейся тестовой функции
	isTestingExists = make(map[string]struct{})
)

func Execute(requestPath, validatePath, testingPath, errPath string) {
	// зачитывание и парсинг файла со структурами в строку
	s := readStruct(requestPath)

	fs := token.NewFileSet()
	// дерево ast
	astFile, err := parser.ParseFile(fs, "", s, parser.ParseComments)
	if err != nil {
		log.Fatalf("ast file parse error: %s", err)
	}

	// проверка наличия папки errors, если нет, то создать
	if err := os.Mkdir(errPath[:len(errPath)-10], os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// проверка на необходимость перезаписи файла validate.go. Перезаписывается шапка.
	if needReWriteFile(validatePath) {
		reWriteFile(validatePath, templates.HeadValidate)
	}
	// проверка на необходимость перезаписи файла errors.go. Перезаписывается шапка.
	if needReWriteFile(errPath) {
		reWriteFile(errPath, templates.HeadErrors)
	}
	// проверка на необходимость перезаписи файла validate_test.go. Перезаписывается шапка.
	if needReWriteFile(testingPath) {
		reWriteFile(testingPath, func() string { return "" })
	}

	// открытие файлов на дозапись
	var (
		fileValidate = fileOpenAppendMode(validatePath)
		fileErr      = fileOpenAppendMode(errPath)
		fileTesting  = fileOpenAppendMode(testingPath)
	)

	var (
		// создание буффера для шаблона "ошибка"
		errBuffer = templates.NewTemplateError()
		// создание буффера для шаблона функции "Validate"
		invokeBuffer = templates.NewInvokeTemplate()
		// создание буффера для шаблона "тестирование"
		testingBuffer = templates.NewTestingTemplate()
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
					// название валидирующей функции | ошибки
					funcName, keyErrTemplate string
				)

				// палучение правил тега
				for _, rules := range rulesParsed {
					// расчет отступа для взятия подстроки для конкатенации буффера
					indent := indentConcat(isFirstConcatFunc, funcBuffer.Buffer)
					// основная логика
					switch {
					// обработка правила "rq"
					case len(rules) == 1 && rules[0] == require:
						// создание шаблона "ошибки"
						errTemplate := templates.RequireErr()
						// проверка на наличие уже созданной ошибки
						keyErrTemplate := require
						if _, ok := isErrExists[keyErrTemplate]; !ok {
							// добавление шаблона ошибки в буфер
							templates.AddErrTemplateToBuffer(isErrExists, keyErrTemplate, errTemplate, isFirstConcatErr, errBuffer)
							isFirstConcatErr = false
						}

						// заполнение полей для шаблона "функции"
						funcFields.FieldsFill(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name)
						// добавление шаблона "функции" в буфер
						funcBuffer.BufferConcat(templates.Require(), indent, isFirstConcatFunc)
						isFirstConcatFunc = false

						// название функции
						funcName = typeSpec.Name.String() + field.Names[0].String()

						// проверка на наличие уже существующий тестирующий функции
						if _, ok := isTestingExists[funcName]; !ok {
							// добавление шаблона "тестирование" в буфер
							templates.AddTestingFuncToBuffer(
								isTestingExists,
								funcName,
								templates.TestingFuncTemplate(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name),
								isFirstConcatTesting,
								testingBuffer)
							isFirstConcatTesting = false
						}

						// проверка на наличие уже вызванной функции в шаблоне "Validate"
						if _, ok := isFuncCall[funcName]; !ok {
							// добавление функции проверки в шаблон "Validate"
							templates.AddFuncCallToBuffer(isFuncCall, funcName, templates.CallingFuncTemplate(funcName), isFirstConcatValidation, invokeBuffer)
							isFirstConcatValidation = false
						}
					// обработка правила "lt"
					case rules[0] == lessThan:
						// создание шаблона "ошибки"
						errTemplate = templates.CreateErrorTemplate(greaterThen, rules[1])
						// получение названия переменной "ошибки"
						errVarName = strings.Split(errTemplate, " ")[5]
						// проверка на наличие уже созданной ошибки
						keyErrTemplate = rules[0] + rules[1]
						if _, ok := isErrExists[keyErrTemplate]; !ok {
							// добавление шаблона ошибки в буфер
							templates.AddErrTemplateToBuffer(isErrExists, keyErrTemplate, errTemplate, isFirstConcatErr, errBuffer)
							isFirstConcatErr = false
						}

						// заполнение полей для шаблона "функции"
						funcFields.FieldsFill(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name)
						funcFields.LessThan = rules[1]
						// создание шаблона "функции"
						funcTemplate := templates.CreateLtFunctionTemplate(isArr, errVarName, rulesParsed)
						// добавление шаблона "функции" в буфер
						funcBuffer.BufferConcat(funcTemplate, indent, isFirstConcatFunc)
						isFirstConcatFunc = false

						// название функции
						funcName = typeSpec.Name.String() + field.Names[0].String()

						// проверка на наличие уже существующий тестирующий функции
						if _, ok := isTestingExists[funcName]; !ok {
							// добавление шаблона "тестирование" в буфер
							templates.AddTestingFuncToBuffer(
								isTestingExists,
								funcName,
								templates.TestingFuncTemplate(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name),
								isFirstConcatTesting,
								testingBuffer)
							isFirstConcatTesting = false
						}

						// проверка на наличие уже вызванной функции в шаблоне "Validate"
						if _, ok := isFuncCall[funcName]; !ok {
							// добавление функции проверки в шаблон "Validate"
							templates.AddFuncCallToBuffer(isFuncCall, funcName, templates.CallingFuncTemplate(funcName), isFirstConcatValidation, invokeBuffer)
							isFirstConcatValidation = false
						}
					// обработка правила "gt"
					case rules[0] == greaterThan:
						// создание шаблона "ошибки"
						errTemplate = templates.CreateErrorTemplate(lessThen, rules[1])
						// получение названия переменной "ошибки"
						errVarName = strings.Split(errTemplate, " ")[5]
						// проверка на наличие уже созданной ошибки
						keyErrTemplate = rules[0] + rules[1]
						if _, ok := isErrExists[keyErrTemplate]; !ok {
							// добавление шаблона ошибки в буфер
							templates.AddErrTemplateToBuffer(isErrExists, keyErrTemplate, errTemplate, isFirstConcatErr, errBuffer)
							isFirstConcatErr = false
						}

						// заполнение полей для шаблона "функции"
						funcFields.FieldsFill(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name)
						funcFields.GreaterThan = rules[1]
						// создание шаблона "функции"
						funcTemplate := templates.CreateGtFunctionTemplate(isArr, errVarName, rulesParsed)
						// добавление шаблона "функции" в буфер
						funcBuffer.BufferConcat(funcTemplate, indent, isFirstConcatFunc)
						isFirstConcatFunc = false

						// название функции
						funcName = typeSpec.Name.String() + field.Names[0].String()

						// проверка на наличие уже существующий тестирующий функции
						if _, ok := isTestingExists[funcName]; !ok {
							// добавление шаблона "тестирование" в буфер
							templates.AddTestingFuncToBuffer(
								isTestingExists,
								funcName,
								templates.TestingFuncTemplate(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name),
								isFirstConcatTesting,
								testingBuffer)
							isFirstConcatTesting = false
						}

						// проверка на наличие уже вызванной функции в шаблоне "Validate"
						if _, ok := isFuncCall[funcName]; !ok {
							// добавление функции проверки в шаблон "Validate"
							templates.AddFuncCallToBuffer(isFuncCall, funcName, templates.CallingFuncTemplate(funcName), isFirstConcatValidation, invokeBuffer)
							isFirstConcatValidation = false
						}
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
