// Package errors - пакет с ошибками стратегии.
package errors

import "errors"

// Is - check wrap errors.
var Is = errors.Is

var (
	// ErrLessThen1 - поле меньше 1 символов.
	ErrLessThen1 = errors.New("field is less than 1 char")
	// ErrGreaterThen30 - поле больше 30 символов.
	ErrGreaterThen30 = errors.New("field is greater than 30 char")
	// ErrLessThen2 - поле меньше 2 символов.
	ErrLessThen2 = errors.New("field is less than 2 char")
	// ErrGreaterThen20 - поле больше 20 символов.
	ErrGreaterThen20 = errors.New("field is greater than 20 char")
	// ErrGreaterThen10 - поле больше 10 символов.
	ErrGreaterThen10 = errors.New("field is greater than 10 char")
	// ErrRequired - обязательное поле.
	ErrRequired = errors.New("field is required")
)
