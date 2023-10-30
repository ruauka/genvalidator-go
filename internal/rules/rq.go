package rules

import (
	"go/ast"

	"github.com/fatih/structtag"

	"genvalidator/internal/templates"
)

var falseStub = false

func Require(
	isFirstConcatFunc, isFirstConcatErr, isFirstConcatValidation, isFirstConcatTesting *bool,
	isErrExists, isTestingExists, isFuncCall map[string]struct{},
	errBuffer *templates.ErrorTemplate, funcBuffer *templates.FunctionTemplate,
	invokeBuffer *templates.InvokeTemplate, testingBuffer *templates.TestingTemplate,
	funcFields *templates.TemplateFields,
	typeSpec *ast.TypeSpec, field *ast.Field, jsonFieldName *structtag.Tag,
	indent int,
) {
	// создание шаблона "ошибки"
	errTemplate := templates.RequireErr()
	// проверка на наличие уже созданной ошибки
	keyErrTemplate := "rq"
	if _, ok := isErrExists[keyErrTemplate]; !ok {
		// добавление шаблона ошибки в буфер
		templates.AddErrTemplateToBuffer(isErrExists, keyErrTemplate, errTemplate, *isFirstConcatErr, errBuffer)
		*isFirstConcatErr = false
	}

	// заполнение полей для шаблона "функции"
	funcFields.FieldsFill(typeSpec.Name.String(), field.Names[0].String(), jsonFieldName.Name)
	// добавление шаблона "функции" в буфер
	funcBuffer.BufferConcat(templates.Require(), indent, *isFirstConcatFunc)
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
