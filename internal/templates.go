package internal

import (
	"fmt"
	"strings"
)

type Template struct {
	Buffer string
}

func NewTemplate() *Template {
	return &Template{Buffer: ValidateFunc()}
}

func (r *Template) Concat(template string, index int, isFirstConcat bool) {
	if isFirstConcat {
		r.Buffer = r.Buffer[:index] + template + r.Buffer[index:]
		return
	}

	r.Buffer = "\n" + strings.TrimSpace(r.Buffer[:index]) + template + "\n" + r.Buffer[index:]
}

func (r *Template) Reset() {
	r.Buffer = ValidateFunc()
}

type TemplateFields struct {
	StructName        string
	StrategyFieldName string
	JsonFieldName     string
	LessThan          string
	GreaterThan       string
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
