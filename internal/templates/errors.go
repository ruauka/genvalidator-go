package templates

import (
	"fmt"
	"strings"
)

const (
	errIndex = 6
)

var (
	errTemplates = make(map[string]struct{})
)

type TemplateError struct {
	Buffer string
}

func NewTemplateError() *TemplateError {
	return &TemplateError{Buffer: Vars()}
}

func (t *TemplateError) Concat(template string, index int, counter int) {
	if counter == 1 {
		t.Buffer = t.Buffer[:index] + template + t.Buffer[index+1:]
		return
	}
	if counter == 2 {
		t.Buffer = strings.TrimSpace(t.Buffer[:len(t.Buffer)-index]) + template + strings.TrimSpace(t.Buffer[len(t.Buffer)-index+1:])
		return
	}

	t.Buffer = strings.TrimSpace(t.Buffer[:len(t.Buffer)-index]) + template + strings.TrimSpace(t.Buffer[len(t.Buffer)-index-1:])
}

// CreateErrorTemplate - создание шаблона ошибки.
func CreateErrorTemplate(prefix []string, val string) string {
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

func AddErrTemplateToBuffer(key, errTemplate string, counter *int, t *TemplateError) {
	// проверка на наличие уже созданной ошибки
	if _, ok := errTemplates[key]; !ok {
		// создание шаблона error
		t.Concat(errTemplate, errIndex, *counter)
		errTemplates[key] = struct{}{}
		*counter++
	}
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
