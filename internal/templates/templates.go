package templates

import "strings"

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
	GreaterThanOrEq   string
}

func Head() string {
	return `package main

import "fmt"

func Len(text string) int {
    return len([]rune(text))
}
`
}

func ValidateFunc() string {
	return `
func {{.StructName}}{{.StrategyFieldName}}(req *Request) error {
    return nil
}
`
}

func Require() string {
	return `
    // require (rq)
    if req.{{.StructName}}.{{.StrategyFieldName}} == nil {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %s", "field is require")
    }
`
}

func LessThan() string {
	return `
    // less than (lt)
    if Len(req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %s", "field must be less than {{.LessThan}}")
    }
`
}

func LessThanPtr() string {
	return `
    // less than (lt)
    if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) > {{.LessThan}} {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %s", "field must be less than {{.LessThan}}")
    }
`
}

func GreaterThanOrEq() string {
	return `
    // greater than or eq (gte)
    if Len(req.{{.StructName}}.{{.StrategyFieldName}}) <= {{.GreaterThanOrEq}} {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %s", "field must be greater than or eq {{.GreaterThanOrEq}}")
    }
`
}

func GreaterThanOrEqPtr() string {
	return `
    // greater than or eq (gte)
    if Len(*req.{{.StructName}}.{{.StrategyFieldName}}) <= {{.GreaterThanOrEq}} {
        return fmt.Errorf("failed check field '{{.StructName}}.{{.JsonFieldName}}': %s", "field must be greater than or eq {{.GreaterThanOrEq}}")
    }
`
}
