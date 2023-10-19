package internal

import (
	"fmt"
	"strings"
)

type TemplateError struct {
	Buffer string
}

func NewTemplateError() *TemplateError {
	return &TemplateError{Buffer: Vars()}
}

func (r *TemplateError) Concat(template string, index int, counter int) {
	if counter == 1 {
		r.Buffer = r.Buffer[:index] + template + r.Buffer[index+1:]
		return
	}
	if counter == 2 {
		r.Buffer = strings.TrimSpace(r.Buffer[:len(r.Buffer)-index]) + template + strings.TrimSpace(r.Buffer[len(r.Buffer)-index+1:])
		return
	}

	r.Buffer = strings.TrimSpace(r.Buffer[:len(r.Buffer)-index]) + template + strings.TrimSpace(r.Buffer[len(r.Buffer)-index-1:])
}

func createErr(prefix []string, val string) string {
	return fmt.Sprintf(
		"\n    // Err%s%s - поле %s %s символов.\n    Err%s%s = errors.New(\"field is %s than %s char\")\n",
		prefix[0], val,
		prefix[1],
		val,
		prefix[0], val,
		prefix[2],
		val,
	)
}

func HeadErrors() string {
	return `package validation

import "errors"
`
}

func Vars() string {
	return `
var (
)
`
}

func RequireErr() string {
	return `
    // ErrRequired - обязательное поле.
    ErrRequired = errors.New("field is required")
`
}

func ErrGreaterThen10() string {
	return `
    // ErrGreaterThen10 - поле больше 10 символов.
    ErrGreaterThen10 = errors.New("field is greater than 10 char")
`
}
