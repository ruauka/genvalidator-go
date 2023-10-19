package validation

import "fmt"

func Len(text string) int {
	return len([]rune(text))
}

func FooField2(req *Request) error {
	// require (rq)
	if req.Foo.Field2 == nil {
		return fmt.Errorf("failed check field 'Foo.field_2': %w", ErrRequired)
	}

	return nil
}

func FooBar(req *Request) error {
	// require (rq)
	if req.Foo.Bar == nil {
		return fmt.Errorf("failed check field 'Foo.bar': %w", ErrRequired)
	}

	return nil
}

func BarField3(req *Request) error {
	// less than (lt)
	if Len(req.Bar.Field3) > 10 {
		return fmt.Errorf("failed check field 'Bar.field_3': %w", ErrGreaterThen10)
	}

	return nil
}

func BarField4(req *Request) error {
	// require (rq)
	if req.Bar.Field4 == nil {
		return fmt.Errorf("failed check field 'Bar.field_4': %w", ErrRequired)
	}
	// less than (lt)
	if Len(*req.Bar.Field4) > 20 {
		return fmt.Errorf("failed check field 'Bar.field_4': %w", ErrGreaterThen20)
	}

	return nil
}

func BazField5(req *Request) error {
	// greater than (gt)
	if Len(req.Baz.Field5) < 2 {
		return fmt.Errorf("failed check field 'Baz.field_5': %w", ErrLessThen2)
	}

	return nil
}

func BazField6(req *Request) error {
	// require (rq)
	if req.Baz.Field6 == nil {
		return fmt.Errorf("failed check field 'Baz.field_6': %w", ErrRequired)
	}
	// greater than (gt)
	if Len(*req.Baz.Field6) < 2 {
		return fmt.Errorf("failed check field 'Baz.field_6': %w", ErrLessThen2)
	}
	// less than (lt)
	if Len(*req.Baz.Field6) > 30 {
		return fmt.Errorf("failed check field 'Baz.field_6': %w", ErrGreaterThen30)
	}

	return nil
}

func BazField7(req *Request) error {
	// require (rq)
	if req.Baz.Field7 == nil {
		return fmt.Errorf("failed check field 'Baz.field_7': %w", ErrRequired)
	}
	// greater than (gt)
	for idx, val := range req.Baz.Field7 {
		if Len(val) < 1 {
			return fmt.Errorf("failed check field 'Baz.field_7': %w, err in %dth array index", ErrLessThen1, idx)
		}
	}
	// less than (lt)
	for idx, val := range req.Baz.Field7 {
		if Len(val) > 10 {
			return fmt.Errorf("failed check field 'Baz.field_7': %w, err in %dth array index", ErrGreaterThen10, idx)
		}
	}

	return nil
}