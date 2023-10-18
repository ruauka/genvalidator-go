package main

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
	Field5 string  `json:"field_5" validate:"gte=2"`
	Field6 *string `json:"field_6" validate:"rq, gte=2, lt=30"`
}
