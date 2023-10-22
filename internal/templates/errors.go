package templates

import (
	"fmt"
	"strings"
)

// отступ для конкатенации
const (
	errIndent = 6
)

// ErrorTemplate - структура шаблона "ошибки".
type ErrorTemplate struct {
	Buffer string
}

// NewTemplateError - конструктор с пердзаполненным шаблоном.
func NewTemplateError() *ErrorTemplate {
	return &ErrorTemplate{Buffer: Vars()}
}

// BufferConcat - добавление шаблона в буффер (конкатенация).
func (e *ErrorTemplate) BufferConcat(template string, indent int, isFirstConcat bool) {
	// первая конкатенация
	if isFirstConcat {
		e.Buffer = e.Buffer[:indent] + template + e.Buffer[indent+1:]
		return
	}

	// последующие конкатенации
	e.Buffer = strings.TrimSpace(e.Buffer[:len(e.Buffer)-indent+3]) + template + strings.TrimSpace(e.Buffer[len(e.Buffer)-indent+4:])
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

// AddErrTemplateToBuffer - добавление новой "ошибки" в буффер.
func AddErrTemplateToBuffer(isErrExists map[string]struct{}, key, errTemplate string, isFirstConcat bool, t *ErrorTemplate) {
	// создание шаблона error
	t.BufferConcat(errTemplate, errIndent, isFirstConcat)
	isErrExists[key] = struct{}{}
}

// HeadErrors - шаблон начала файла.
func HeadErrors() string {
	return `package validation

import "errors"
`
}

// Vars - объявление "ошибок".
func Vars() string {
	return `
var (
)
`
}

// RequireErr - шаблон "ошибки" обязательности поля.
func RequireErr() string {
	return `
    // ErrRequired - обязательное поле.
    ErrRequired = errors.New("field is required")
`
}
