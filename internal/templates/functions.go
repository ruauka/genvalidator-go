package templates

import (
	"fmt"
	"strings"
)

type FunctionsTemplate struct {
	Buffer string
}

func NewFunctionsTemplate() *FunctionsTemplate {
	return &FunctionsTemplate{Buffer: ValidateFunc()}
}

func (f *FunctionsTemplate) Concat(template string, index int, isFirstConcat bool) {
	if isFirstConcat {
		f.Buffer = f.Buffer[:index] + template + f.Buffer[index:]
		return
	}

	f.Buffer = "\n" + strings.TrimSpace(f.Buffer[:index]) + template + "\n" + f.Buffer[index:]
}

func (f *FunctionsTemplate) Reset() {
	f.Buffer = ValidateFunc()
}

func CreateLtFunctionTemplate(isArr bool, errVarName string, tagsParsed [][]string) string {
	switch {
	case isArr:
		return LessThanSl(errVarName)
	case isRqTag(tagsParsed):
		return LessThanPtr(errVarName)
	default:
		return LessThan(errVarName)
	}
}

func CreateGtFunctionTemplate(isArr bool, errVarName string, tagsParsed [][]string) string {
	switch {
	case isArr:
		return GreaterThanSl(errVarName)
	case isRqTag(tagsParsed):
		return GreaterThanPtr(errVarName)
	default:
		return GreaterThan(errVarName)
	}
}

func isRqTag(tags [][]string) bool {
	for _, tag := range tags {
		if tag[0] == "rq" {
			return true
		}
	}

	return false
}

type TemplateFields struct {
	StructName        string
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

func HeadValidate() string {
	return `package validation

import "fmt"

// Len - длина строки (в рунах).
func Len(text string) int {
    return len([]rune(text))
}
`
}

func ValidateFunc() string {
	return `
// {{.StructName}}{{.StrategyFieldName}} - валидация поля {{.StructName}}.{{.JsonFieldName}}.
func {{.StructName}}{{.StrategyFieldName}}(req *Request) error {
    return nil
}
`
}

func Require() string {
	return `
    // require (rq)
    if req.{{.StructName}}.{{.StrategyFieldName}} == nil {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %w", ErrRequired)
    }
`
}

func LessThan(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"if Len(req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {\n"+
		"        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", %s)\n    }\n", errStr)
}

func LessThanPtr(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {\n"+
		"        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", %s)\n    }\n", errStr)
}

func LessThanSl(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// less than (lt)\n    "+
		"for idx, val := range req.{{.StructName}}.{{.StrategyFieldName}} {"+
		"\n        if Len(val) > {{.LessThan}} {"+
		"\n            return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w, err in %%dth array index\", %s, idx)\n        }\n    }\n", errStr)
}

func GreaterThan(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// greater than (gt)\n    "+
		"if Len(req.{{.StructName}}.{{.StrategyFieldName}}) < {{.GreaterThan}}"+
		" {\n        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", %s)\n    }\n", errStr)
}

func GreaterThanPtr(errStr string) string {
	return fmt.Sprintf("\n    "+
		"// greater than (gt)\n    "+
		"if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) < {{.GreaterThan}}"+
		" {\n        return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w\", %s)\n    }\n", errStr)
}

func GreaterThanSl(errStr string) string {
	return fmt.Sprintf(
		"\n    "+
			"// greater than (gt)\n    "+
			"for idx, val := range req.{{.StructName}}.{{.StrategyFieldName}} {"+
			"\n        if Len(val) < {{.GreaterThan}} {\n"+
			"            return fmt.Errorf(\"failed check field '{{.StructName}}.{{.JsonFieldName}}': %%w, err in %%dth array index\", %s, idx)\n        }\n    }\n", errStr,
	)
}
