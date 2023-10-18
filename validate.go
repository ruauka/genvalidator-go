package main

import "fmt"

func Len(text string) int {
	return len([]rune(text))
}

func FooField2(req *Request) error {
	// require (rq)
	if req.Foo.Field2 == nil {
		return fmt.Errorf("failed check field 'Foo.field_2': %s", "field is require")
	}

	return nil
}

func FooBar(req *Request) error {
	// require (rq)
	if req.Foo.Bar == nil {
		return fmt.Errorf("failed check field 'Foo.bar': %s", "field is require")
	}

	return nil
}

func BarField3(req *Request) error {
	// less than (lt)
	if Len(req.Bar.Field3) > 10 {
		return fmt.Errorf("failed check field 'Bar.field_3': %s", "field must be less than 10")
	}

	return nil
}

func BarField4(req *Request) error {
	// require (rq)
	if req.Bar.Field4 == nil {
		return fmt.Errorf("failed check field 'Bar.field_4': %s", "field is require")
	}
	// less than (lt)
	if Len(*req.Bar.Field4) > 20 {
		return fmt.Errorf("failed check field 'Bar.field_4': %s", "field must be less than 20")
	}

	return nil
}

func BazField5(req *Request) error {
	// greater than or eq (gte)
	if Len(req.Baz.Field5) <= 2 {
		return fmt.Errorf("failed check field 'Baz.field_5': %s", "field must be greater than or eq 2")
	}

	return nil
}

func BazField6(req *Request) error {
	// require (rq)
	if req.Baz.Field6 == nil {
		return fmt.Errorf("failed check field 'Baz.field_6': %s", "field is require")
	}
	// greater than or eq (gte)
	if Len(*req.Baz.Field6) <= 2 {
		return fmt.Errorf("failed check field 'Baz.field_6': %s", "field must be greater than or eq 2")
	}
	// less than (lt)
	if Len(*req.Baz.Field6) > 30 {
		return fmt.Errorf("failed check field 'Baz.field_6': %s", "field must be less than 30")
	}

	return nil
}
