package request_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"

	"genvalidator/validation/errors"
	"genvalidator/validation/request"
)

func TestFooField2(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Foo, field:Field2, json:field_2' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Foo, field:Field2, json:field_2' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("FooField2")
			t.Description("Проверка валидации входящего запроса. Json поле: field_2")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.FooField2(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestFooBar(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Foo, field:Bar, json:bar' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Foo, field:Bar, json:bar' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("FooBar")
			t.Description("Проверка валидации входящего запроса. Json поле: bar")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.FooBar(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestBarField3(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Bar, field:Field3, json:field_3' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Bar, field:Field3, json:field_3' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("BarField3")
			t.Description("Проверка валидации входящего запроса. Json поле: field_3")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.BarField3(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestBarField4(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Bar, field:Field4, json:field_4' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Bar, field:Field4, json:field_4' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("BarField4")
			t.Description("Проверка валидации входящего запроса. Json поле: field_4")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.BarField4(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestBazField5(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field5, json:field_5' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field5, json:field_5' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("BazField5")
			t.Description("Проверка валидации входящего запроса. Json поле: field_5")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.BazField5(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestBazField6(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field6, json:field_6' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field6, json:field_6' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("BazField6")
			t.Description("Проверка валидации входящего запроса. Json поле: field_6")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.BazField6(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}

func TestBazField7(tr *testing.T) {
	// создаем объект test runner
	r := runner.NewRunner(tr, tr.Name())

	testCases := []struct {
		File  string
		Name  string
		Error error
	}{
		{
			Name:  "Успешная валидация",
			File:  "success.json",
			Error: nil,
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field7, json:field_7' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case1.json",
			Error: nil, // case1 error. Взять ошубку из пакета pkg/errors
		},
		{
			Name:  "Ошибка валидации 'struct:Baz, field:Field7, json:field_7' ...описание ошибки: поле обязательное, больше чем столько-то и тд...",
			File:  "case2.json",
			Error: nil, // case2 error. Взять ошубку из пакета pkg/errors
		},
	}

	for _, testCase := range testCases {
		test := testCase

		r.NewTest(test.Name, func(t provider.T) {
			var req request.Request

			// задаем allure id
			t.AllureID(fmt.Sprintf("%s_%s", tr.Name(), test.Name))

			// указываем информацию о тестах для allure отчета
			t.Epic("Validating")
			t.Story("BazField7")
			t.Description("Проверка валидации входящего запроса. Json поле: field_7")

			// чтение JSON файла
			path := filepath.Join("../../testdata/validate", test.File)
			body, err := os.ReadFile(path)
			t.Assert().NoError(err, "Чтение тестовго файла")

			// сериализация JSON в структуру request.Request
			err = json.Unmarshal(body, &req)
			t.Assert().NoError(err, fmt.Sprintf("Преобразование %s", path))

			err = request.BazField7(&req)
			params := allure.NewParameters("Expected", fmt.Sprintf("%v", test.Error), "Actual", fmt.Sprintf("%v", err))

			// временный assert на тип error
			t.WithNewStep("ASSERT: ErrorIs", func(sCtx provider.StepCtx) {
				if !errors.Is(err, test.Error) { // импорт из pkg/errors
					sCtx.FailNow()
				}
			}, params...)
			t.WithAttachments(allure.NewAttachment("Request", allure.JSON, body))
		})
	}

	r.RunTests()
}
