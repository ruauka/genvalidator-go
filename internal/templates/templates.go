package templates

type Data struct {
	StructName string
	FieldName  string
}

func New(structName, fieldName string) *Data {
	return &Data{
		StructName: structName,
		FieldName:  fieldName,
	}
}

func Head() string {
	return `package main

import "fmt"
`
}

func Require() string {
	return `
func {{.StructName}}{{.FieldName}}Check(req *Request) error {
	if req.{{.StructName}}.{{.FieldName}} == nil {
		return fmt.Errorf("failed check field '{{.StructName}}.{{.FieldName}}': %s", "field is require")
	}

	return nil
}
`
}
