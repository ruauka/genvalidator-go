package rules

import (
	"go/ast"
	"strings"

	"github.com/fatih/structtag"

	"genvalidator/internal/templates"
)

var (
	greaterThen = []string{"GreaterThen", "больше", "greater"}
)

func LessThan(
	isFirstConcatFunc, isFirstConcatErr, isFirstConcatValidation, isFirstConcatTesting *bool,
	isErrExists, isTestingExists, isFuncCall map[string]struct{},
	errBuffer *templates.ErrorTemplate, funcBuffer *templates.FunctionTemplate,
	invokeBuffer *templates.InvokeTemplate, testingBuffer *templates.TestingTemplate,
	funcFields *templates.TemplateFields,
	typeSpec *ast.TypeSpec, field *ast.Field, jsonFieldName *structtag.Tag,
	indent int,
	isArr bool,
	rules []string,
	rulesParsed [][]string,
) {
	// создание шаблона "ошибки"
	errTemplate := templates.CreateErrorTemplate(greaterThen, rules[1])
	// получение названия переменной "ошибки"
	errVarName := strings.Split(errTemplate, " ")[5]
	// проверка на наличие уже созданной ошибки
	keyErrTemplate := rules[0] + rules[1]
	if _, ok := isErrExists[keyErrTemplate]; !ok {
		// добавление шаблона ошибки в буфер
		templates.AddErrTemplateToBuffer(isErrExists, keyErrTemplate, errTemplate, *isFirstConcatErr, errBuffer)
		isFirstConcatErr = &falseStub
	}

	// заполнение полей для шаблона "функции"
	funcFields.FieldsFill(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name)
	funcFields.LessThan = rules[1]
	// создание шаблона "функции"
	funcTemplate := templates.CreateLtFunctionTemplate(isArr, errVarName, rulesParsed)
	// добавление шаблона "функции" в буфер
	funcBuffer.BufferConcat(funcTemplate, indent, *isFirstConcatFunc)
	isFirstConcatFunc = &falseStub

	// название функции
	funcName := typeSpec.Name.String() + field.Names[0].String()

	// проверка на наличие уже существующий тестирующий функции
	if _, ok := isTestingExists[funcName]; !ok {
		// добавление шаблона "тестирование" в буфер
		templates.AddTestingFuncToBuffer(
			isTestingExists,
			funcName,
			templates.TestingFuncTemplate(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name),
			*isFirstConcatTesting,
			testingBuffer,
		)
		isFirstConcatTesting = &falseStub
	}

	// проверка на наличие уже вызванной функции в шаблоне "Validate"
	if _, ok := isFuncCall[funcName]; !ok {
		// добавление функции проверки в шаблон "Validate"
		templates.AddFuncCallToBuffer(isFuncCall, funcName, templates.CallingFuncTemplate(funcName), *isFirstConcatValidation, invokeBuffer)
		isFirstConcatValidation = &falseStub
	}
}
