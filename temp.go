package main

type Request struct {
	Foo
	Bar
}

type Foo struct {
	Field1 string   `json:"field_1"`
	Field2 *float64 `json:"field_2" validate:"rq"`
	Bar    []Bar
}

type Bar struct {
	Field3 *string `json:"field_3" validate:"rq"`
}

//func FooField2Check(req *Request) error {
//	if req.Foo.Field2 == nil {
//		return fmt.Errorf("failed check field 'OCONTROL': %s", "field is require")
//	}
//
//	return nil
//}
