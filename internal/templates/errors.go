package templates

import (
	"fmt"
	"strings"
)

// отступ для конкатенации
const (
	errIndent = 6
)

// мапа для проверки уже имеющийся "ошибки"
var (
	isErrExists = make(map[string]struct{})
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
func (t *ErrorTemplate) BufferConcat(template string, indent int, counter int) {
	// первая конкатенация
	if counter == 1 {
		t.Buffer = t.Buffer[:indent] + template + t.Buffer[indent+1:]
		return
	}
	// вторая конкатенация
	if counter == 2 {
		t.Buffer = strings.TrimSpace(t.Buffer[:len(t.Buffer)-indent]) + template + strings.TrimSpace(t.Buffer[len(t.Buffer)-indent+1:])
		return
	}
	// последующие конкатенации
	t.Buffer = strings.TrimSpace(t.Buffer[:len(t.Buffer)-indent]) + template + strings.TrimSpace(t.Buffer[len(t.Buffer)-indent-1:])
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

// AddErrTemplateToBuffer - проверка наличия уже созданной "ошибки" и добавление новой в буффер.
func AddErrTemplateToBuffer(key, errTemplate string, counter *int, t *ErrorTemplate) {
	// проверка на наличие уже созданной ошибки
	if _, ok := isErrExists[key]; !ok {
		// создание шаблона error
		t.BufferConcat(errTemplate, errIndent, *counter)
		isErrExists[key] = struct{}{}
		*counter++
	}
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
