package merger

import (
	"reflect"
	"testing"
)

type test struct {
}

func (t *test) Merge(dest reflect.Value, src reflect.Value) error {
	return nil
}

type Test struct {
	Source interface{}
	Dest   interface{}

	Expect interface{}
}

type EmbeddedStructA struct {
	C string
}

type TestStructA struct {
	EmbeddedStructA

	A string
	B string
}

type TestStructC struct {
	A string
	B string
	C string
}

var tests = []Test{
	Test{
		Source: struct {
			A string
			B string
		}{
			A: "test",
			B: "test",
		},
		Dest: &struct {
			A string
			B string
		}{
			A: "overwritten",
			B: "overwritten",
		},
		Expect: &struct {
			A string
			B string
		}{
			A: "test",
			B: "test",
		},
	},
	Test{
		Source: TestStructA{
			EmbeddedStructA: EmbeddedStructA{
				C: "from-embedded-c",
			},
			A: "test",
			B: "test",
		},
		Dest: &TestStructC{
			A: "overwritten",
			B: "overwritten",
			C: "overwritten",
		},
		Expect: &struct {
			A string
			B string
			C string
		}{
			A: "test",
			B: "test",
			C: "from-embedded-c",
		},
	},
}

func TestXxx(t *testing.T) {
	for _, test := range tests {
		if err := Merge(&test.Dest, test.Source); err != nil {
			t.Errorf("Merge failed: %s", err.Error())
		}

		if reflect.DeepEqual(test.Dest, test.Expect) {
			continue
		}

		t.Errorf("Object not what expected. Got %#v, expected %#v.", test.Dest, test.Expect)
	}
}
