## GenValidator

## Overview

Validating code generator.

## Content

- [Usage](#usage)
- [Rules](#rules)
- [Result](#result)

## Usage

1. Download binary file `genvalidate` and put in project root:
```bash
curl -OL https://raw.githubusercontent.com/ruauka/genvalidator-go/master/genvalidate && chmod +x genvalidate
```
2. Create file `request.go` with your struct in folder `pkg -> request`. 

```bash
.
├── genvalidate
├── go.mod
├── go.sum
├── main.go
└── pkg
    └── request
        └── request.go
```

```go
package request

type Request struct {
    Foo `json:"foo"`
    Bar `json:"bar"`
    Baz `json:"baz"`
}

type Foo struct {
    Field1 string   `json:"field_1"`
    Field2 *float64 `json:"field_2" validate:"rq"`
    Bar    []Bar    `json:"bar" validate:"rq"`
}

type Bar struct {
    Field3 string  `json:"field_3" validate:"lt=10"`
    Field4 *string `json:"field_4" validate:"rq, lt=20"`
}

type Baz struct {
    Field5 string   `json:"field_5" validate:"gt=2"`
    Field6 *string  `json:"field_6" validate:"rq, gt=2, lt=30"`
    Field7 []string `json:"field_7" validate:"rq, gt=1, lt=10"`
}
```

3. Use tag `validate` and 3 rules.
4. Add string `//go:generate ../../genvalidate` above your Request struct.

```go
package request

//go:generate ../../genvalidate

type Request struct {
    Foo `json:"foo"`
    Bar `json:"bar"`
    Baz `json:"baz"`
}
...

```

5. Run script:

```bash
go generate ./...
```

## Rules

- rq - Required field. All field types.
- lt - Length of the field must be less than the specified length. Only `string` and slices with string type.
- gt - Length of the field must be greater than the specified length. Only `string` and slices with string type.

## Result
Script generate 1 folder(if not exists) and 3 files:
    
- `errors.go` (package errors)
- `validate.go`, `validate_test.go` (package request)

```bash
.
├── genvalidate
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── errors
    │   └── errors.go
    └── request
        ├── request.go
        ├── validate.go
        └── validate_test.go
```