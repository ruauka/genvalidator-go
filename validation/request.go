package validation

type Request struct {
	Foo
	Bar
	Baz
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
