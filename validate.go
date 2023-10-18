package main

import "fmt"

func FooField2Check(req *Request) error {
	if req.Foo.Field2 == nil {
		return fmt.Errorf("failed check field 'Foo.Field2': %s", "field is require")
	}

	return nil
}

func BarField3Check(req *Request) error {
	if req.Bar.Field3 == nil {
		return fmt.Errorf("failed check field 'Bar.Field3': %s", "field is require")
	}

	return nil
}
