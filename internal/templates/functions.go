package templates

import (
	"fmt"
	"strings"
)

// FunctionTemplate - структура шаблона "функции" валлидации.
type FunctionTemplate struct {
	Buffer string
}

// NewFunctionsTemplate - конструктор с пердзаполненным шаблоном.
func NewFunctionsTemplate() *FunctionTemplate {
	return &FunctionTemplate{Buffer: RuleFuncName()}
}

// BufferConcat - добавление шаблона в буффер (конкатенация).
func (f *FunctionTemplate) BufferConcat(template string, indent int, isFirstConcat bool) {
	// проверка на первую конкатенацию буфера
	if isFirstConcat {
		f.Buffer = f.Buffer[:indent] + template + f.Buffer[indent:]
		return
	}

	f.Buffer = "\n" + strings.TrimSpace(f.Buffer[:indent]) + template + "\n" + f.Buffer[indent:]
}

// Reset - сброс буфера для создания новой функции валидации поля. Предзаполнение шапкой функции.
func (f *FunctionTemplate) Reset() {
	f.Buffer = RuleFuncName()
}

// CreateLtFunctionTemplate - выбор шаблона lt.
func CreateLtFunctionTemplate(isArr bool, errVarName string, tagsParsed [][]string) string {
	switch {
	case isArr:
		return LessThanSl(errVarName)
	case isRqRule(tagsParsed):
		return LessThanPtr(errVarName)
	default:
		return LessThan(errVarName)
	}
}

// CreateGtFunctionTemplate - выбор шаблона gt.
func CreateGtFunctionTemplate(isArr bool, errVarName string, tagsParsed [][]string) string {
	switch {
	case isArr:
		return GreaterThanSl(errVarName)
	case isRqRule(tagsParsed):
		return GreaterThanPtr(errVarName)
	default:
		return GreaterThan(errVarName)
	}
}

// isRqRule - проверка на rq.
func isRqRule(tags [][]string) bool {
	for _, tag := range tags {
		if tag[0] == "rq" {
			return true
		}
	}

	return false
}

// TemplateFields - структура полей для подстановки в шаблон "функции"
type TemplateFields struct {
	StructName        string
	FuncName          string
	StrategyFieldName string
	JsonFieldName     string
	LessThan          string
	GreaterThan       string
}

// FieldsFill - заполнение полей для шаблона "функции".
func (t *TemplateFields) FieldsFill(structN, fieldN, jsonN string) {
	t.StructName = structN
	t.StrategyFieldName = fieldN
	t.JsonFieldName = jsonN
}

// HeadValidate - шаблон начала файла.
func HeadValidate() string {
	return `package request

import "fmt"

// Len - длина строки (в рунах).
func Len(text string) int {
    return len([]rune(text))
}
`
}

// RuleFuncName - шапка "функции".
func RuleFuncName() string {
	return `
// {{.StructName}}{{.StrategyFieldName}} - валидация поля {{.StructName}}.{{.JsonFieldName}}.
func {{.StructName}}{{.StrategyFieldName}}(req *Request) error {
    return nil
}
`
}

// Require - шаблон обязательности поля.
func Require() string {
	return `
    // require (rq)
    if req.{{.StructName}}.{{.StrategyFieldName}} == nil {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %w", errors.ErrRequired)
    }
`
}

// LessThan - шаблон проверки на "меньше чем" - lt.
func LessThan(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"if Len(req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {\n"+
		"        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", errors.%s)\n    }\n", errStr)
}

// LessThanPtr - шаблон проверки на "меньше чем" и поле обязательное - lt.
func LessThanPtr(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {\n"+
		"        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", errors.%s)\n    }\n", errStr)
}

// LessThanSl - шаблон проверки элементов слайса строк на "меньше чем" - lt.
func LessThanSl(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"for idx, val := range req.{{.StructName}}.{{.StrategyFieldName}} {"+
		"\n        if Len(val) > {{.LessThan}} {"+
		"\n            return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w, err in %%dth array index\", errors.%s, idx)\n        }\n    }\n", errStr)
}

// GreaterThan - шаблон проверки на "больше чем" - gt.
func GreaterThan(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// greater than (gt)\n    "+
		"if Len(req.{{.StructName}}.{{.StrategyFieldName}}) < {{.GreaterThan}}"+
		" {\n        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", errors.%s)\n    }\n", errStr)
}

// GreaterThanPtr - шаблон проверки на "больше чем" и поле обязательное - gt.
func GreaterThanPtr(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// greater than (gt)\n    "+
		"if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) < {{.GreaterThan}}"+
		" {\n        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", errors.%s)\n    }\n", errStr)
}

// GreaterThanSl - шаблон проверки элементов слайса строк на "больше чем" - gt.
func GreaterThanSl(errStr string) string {
	return fmt.Sprintf(
		"\n    "+
			"// greater than (gt)\n    "+
			"for idx, val := range req.{{.StructName}}.{{.StrategyFieldName}} {"+
			"\n        if Len(val) < {{.GreaterThan}} {\n"+
			"            return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w, err in %%dth array index\", errors.%s, idx)\n        }\n    }\n", errStr,
	)
}

//// Validate - валидация входящего запроса.
//func Validate(req *Request) error {
//	if err := FooField2(req); err != nil {
//		return fmt.Errorf("validate error: %w", err)
//	}
//	if err := FooBar(req); err != nil {
//		return fmt.Errorf("validate error: %w", err)
//	}
//
//	return nil
//}
