package templates

import (
	"fmt"
	"strings"
)

// отступ для конкатенации
const indentTesting = 269

// TestingTemplate - структура шаблона "тестирование".
type TestingTemplate struct {
	Buffer string
}

// NewTestingTemplate - конструктор с пердзаполненным шаблоном.
func NewTestingTemplate() *TestingTemplate {
	return &TestingTemplate{Buffer: TestingHeadTemplate()}
}

// BufferConcat - добавление шаблона в буффер (конкатенация).
func (t *TestingTemplate) BufferConcat(template string, isFirstConcat bool) {
	// первая конкатенация
	if isFirstConcat {
		t.Buffer = t.Buffer[:indentTesting] + "\n" + template + t.Buffer[indentTesting+1:]
		return
	}

	// последующие конкатенации
	t.Buffer = strings.TrimSpace(t.Buffer) + "\n" + template
}

// AddTestingFuncToBuffer - добавление нового вызова новой тестурующий функции в буффер.
func AddTestingFuncToBuffer(isTestingExists map[string]struct{}, key, template string, isFirstConcat bool, t *TestingTemplate) {
	// добавление в буффер
	t.BufferConcat(template, isFirstConcat)
	isTestingExists[key] = struct{}{}
}

// TestingHeadTemplate - шаблон "тестирования".
func TestingHeadTemplate() string {
	return `package request_test

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "testing"

    "github.com/ozontech/allure-go/pkg/allure"
    "github.com/ozontech/allure-go/pkg/framework/provider"
    "github.com/ozontech/allure-go/pkg/framework/runner"
)
`
}

// TestingFuncTemplate - формирование шаблона вызова валидирующий функции.
func TestingFuncTemplate(structName, fieldName, jsonName string) string {
	funcName := structName + fieldName
	return fmt.Sprintf(
		"\nfunc Test%s(t *testing.T) {"+
			"\n    testCases := []struct {\n        "+
			"File  string\n        "+
			"Name  string\n        "+
			"Error error\n    }"+
			"{\n        "+
			"{\n            "+
			"Name:  \"Успешная валидация\",\n            "+
			"File:  \"success.json\",\n            "+
			"Error: nil,\n        "+
			"},\n        {\n            "+
			"Name:  \"Ошибка валидации 'struct:%s, field:%s, json:%s' ...описание ошибки: поле обязательное, больше чем столько-то и тд...\",\n            "+
			"File:  \"case1.json\",\n            "+
			"Error: nil, // case1 error. Взять ошубку из пакета pkg/errors\n        "+
			"},\n        {\n            "+
			"Name:  \"Ошибка валидации 'struct:%s, field:%s, json:%s' ...описание ошибки: поле обязательное, больше чем столько-то и тд...\",\n            "+
			"File:  \"case2.json\",\n            "+
			"Error: nil, // case2 error. Взять ошубку из пакета pkg/errors\n        "+
			"},\n    }\n\n    "+
			"for _, testCase := range testCases {"+
			"\n        test := testCase\n\n        "+
			"runner.Run(t, test.Name, func(t provider.T) {"+
			"\n            var req request.Request\n\n            "+
			"// задаем allure id\n            "+
			"t.AllureID(t.Name())\n\n"+
			"            // указываем информацию о тестах для allure отчета\n            "+
			"t.Epic(\"Validating\")\n            "+
			"t.Story(\"%s\")\n            "+
			"t.Description(\"Проверка валидации входящего запроса. Json поле: %s\")\n\n            "+
			"// чтение JSON файла\n            "+
			"path := filepath.Join(\"../../testdata/validate\", "+"test.File)\n           "+
			" body, err := os.ReadFile(path)\n            "+
			"t.Assert().NoError(err, \"Чтение тестовго файла\")\n\n            "+
			"// сериализация JSON в структуру request.Request\n            "+
			"err = json.Unmarshal(body, &req)\n            "+
			"t.Assert().NoError(err, fmt.Sprintf(\"Преобразование %%s\", path))\n\n            "+
			"err = request.%s(&req)\n            "+
			"params := allure.NewParameters(\"Expected\", fmt.Sprintf(\"%%v\", test.Error), \"Actual\", fmt.Sprintf(\"%%v\", err))\n\n"+
			"            // временный assert на тип error\n            "+
			"t.WithNewStep(\"ASSERT: ErrorIs\", func(sCtx provider.StepCtx) {"+
			"\n                if !errors.Is(err, test.Error) { // импорт из pkg/errors\n                   "+
			" sCtx.FailNow()\n                }\n            }, params...)\n            "+
			"t.WithAttachments(allure.NewAttachment(\"Request\", allure.JSON, body))\n        })\n    }\n"+
			"}\n", funcName, structName, fieldName, jsonName, structName, fieldName, jsonName, funcName, jsonName, funcName,
	)
}
