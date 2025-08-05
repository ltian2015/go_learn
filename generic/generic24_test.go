package generic

import (
	"fmt"
	"testing"
)

type GenericStruct[T1, T2 any] struct {
	id    T1
	Value T2
}
type IntIdStruct[T any] struct {
	id    int
	Value T
}

func TestGenericStruct(t *testing.T) {
	var gsint = GenericStruct[int, int]{
		id:    1,
		Value: 100,
	}
	var gsString = GenericStruct[int, string]{
		id:    2,
		Value: "hello",
	}

	var gsintIdStrValue = IntIdStruct[string]{
		id:    3,
		Value: "world",
	}
	var gsintIdIntValue = IntIdStruct[int]{
		id:    20,
		Value: 200,
	}
	println(gsint.id)
	println(gsString.Value)
	printContent(gsint)
	printContent(gsString)
	idAdd(gsintIdStrValue)
	idAdd(gsintIdIntValue)

}
func idAdd[T any](iis IntIdStruct[T]) {
	iis.id += 2000
	fmt.Printf("%v+10000=%v\n", iis.id, iis.id+10000)
}
func printContent[T1, T2 any](gs GenericStruct[T1, T2]) {
	fmt.Printf("{id: %v, Value: %v }\n", gs.id, gs.Value)
}

type structField interface {
	struct {
		a int
		X int
	} |
		struct {
			b int
			X int
		} |
		struct {
			c int
			X int
		}
}

// This function is INVALID.
func IncrementX2[T structField](p T) {
	//v := p.X // INVALID: type of p.x is not the same for all types in set
	//v++
	//	p.X = v
}
