## GenValidator

## Overview

Validating code generator.

## Rules

- rq - Required field. All field types.
- lt - Length of the field must be less than the specified length. Only `string` and slices with string type.
- gt - Length of the field must be greater than the specified length. Only `string` and slices with string type.

## Usage

1. Add binary file `genvalidate` in project root.
2. Create file `request.go` with your `struct`.

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
4. Add string `//go:generate ./genvalidate pkg/request pkg/errors` above your Execute() func.

```go
package main

import (
	"fmt"
	
	"project/pkg/request"
)

//go:generate ./genvalidate pkg/request pkg/errors

// Execute - main project func.
func Execute(req *request.Request) error {
    if err := request.Validate(req); err != nil {
        return fmt.Errorf("validation err: %w", err)
    }
    
    return nil
}

func main() {}

```

- pkg/request - path where script creates `validate.go`, `validate_test.go`
- pkg/errors - path where script creates `errors.go`

5. Run plugin:

```bash
go generate ./...
```