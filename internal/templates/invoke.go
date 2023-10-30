package templates

import (
	"fmt"
	"strings"
)

// InvokeTemplate - структура шаблона "функции validate".
type InvokeTemplate struct {
	Buffer string
}

// NewInvokeTemplate - конструктор с пердзаполненным шаблоном.
func NewInvokeTemplate() *InvokeTemplate {
	return &InvokeTemplate{Buffer: ValidateFuncTemplate()}
}

// BufferConcat - добавление шаблона в буффер (конкатенация).
func (i *InvokeTemplate) BufferConcat(template string, isFirstConcat bool) {
	// первая конкатенация
	if isFirstConcat {
		i.Buffer = i.Buffer[:121] + template + i.Buffer[122:]
		return
	}

	// последующие конкатенации
	i.Buffer = strings.TrimSpace(i.Buffer[:len(i.Buffer)-13]) + template + i.Buffer[len(i.Buffer)-18:]
}

// AddFuncCallToBuffer - добавление нового вызова новой валидирующий функции в буффер.
func AddFuncCallToBuffer(isFuncCall map[string]struct{}, key, template string, isFirstConcat bool, t *InvokeTemplate) {
	// добавление в буффер
	t.BufferConcat(template, isFirstConcat)
	isFuncCall[key] = struct{}{}
}

// ValidateFuncTemplate - шаблон функции "Validate".
func ValidateFuncTemplate() string {
	return `
// Validate - валидация входящего запроса.
func Validate(req *Request) error {
	var err error	


    return nil
}
`
}

// CallingFuncTemplate - формирование шаблона вызова валидирующий функции.
func CallingFuncTemplate(funcName string) string {
	return fmt.Sprintf(
		"\n    if err = %s(req); err != nil {"+
			"\n        return fmt.Errorf(\"validate error: %%w\", err)\n    }\n", funcName,
	)
}
