# structtag [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/fatih/structtag) 

structtag provides an easy way of parsing and manipulating struct rules fields.
Please vendor the library as it might change in future versions.

# Install

```bash
go get github.com/fatih/structtag
```

# Example

```go
package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/fatih/structtag"
)

func main() {
	type t struct {
		t string `json:"foo,omitempty,string" xml:"foo"`
	}

	// get field rules
	rules := reflect.TypeOf(t{}).Field(0).Tag

	// ... and start using structtag by parsing the rules
	tags, err := structtag.Parse(string(rules))
	if err != nil {
		panic(err)
	}

	// iterate over all tags
	for _, t := range tags.Tags() {
		fmt.Printf("rules: %+v\n", t)
	}

	// get a single rules
	jsonTag, err := tags.Get("json")
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonTag)         // Output: json:"foo,omitempty,string"
	fmt.Println(jsonTag.Key)     // Output: json
	fmt.Println(jsonTag.Name)    // Output: foo
	fmt.Println(jsonTag.Options) // Output: [omitempty string]

	// change existing rules
	jsonTag.Name = "foo_bar"
	jsonTag.Options = nil
	tags.Set(jsonTag)

	// add new rules
	tags.Set(&structtag.Tag{
		Key:     "hcl",
		Name:    "foo",
		Options: []string{"squash"},
	})

	// print the tags
	fmt.Println(tags) // Output: json:"foo_bar" xml:"foo" hcl:"foo,squash"

	// sort tags according to keys
	sort.Sort(tags)
	fmt.Println(tags) // Output: hcl:"foo,squash" json:"foo_bar" xml:"foo"
}
```
