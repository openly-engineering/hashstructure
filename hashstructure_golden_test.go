package hashstructure

import (
	"bytes"
	"fmt"
	"testing"
)

type goldenStruct struct {
	UUID    string `hash:"ignore"`
	AMap    map[string]structB
	ASlice  []int
	AString string
	APtr    *structB
}

type structB struct {
	A int
	B bool
}

var (
	goldenStructA = goldenStruct{
		UUID: "Foobar",
		AMap: map[string]structB{
			"bar": {
				A: 2,
				B: true,
			},
			"baz": {
				A: -5,
			},
			"bat": {},
		},
		ASlice:  []int{5, 42, -13, 0, 0, 24},
		AString: "hello",
		APtr: &structB{
			A: 4,
			B: true,
		},
	}

	goldenStructB = goldenStruct{
		UUID: "DIFFERENT_THAN_ABOVE",
		AMap: map[string]structB{
			"bar": {
				A: 2,
				B: true,
			},
			"bat": {},
			"baz": {
				A: -5,
			},
		},
		ASlice:  []int{5, 42, -13, 0, 0, 24},
		AString: "hello",
		APtr: &structB{
			A: 4,
			B: true,
		},
	}

	goldenStructC = goldenStruct{}

	goldenStructD = goldenStruct{
		UUID: "uuid",
		AMap: map[string]structB{
			"bar": {
				A: 2,
				B: true,
			},
			"bat": {},
			"baz": {
				A: -5,
			},
			"1": {},
			"2": {},
		},
		ASlice:  nil,
		AString: "",
		APtr:    nil,
	}

	goldenStructs = []goldenStruct{goldenStructA, goldenStructB, goldenStructC, goldenStructD}
	goldenHashes  = [][]byte{
		{115, 205, 154, 57, 182, 201, 82, 233, 17, 152, 239, 179, 145, 124, 147, 33},
		{115, 205, 154, 57, 182, 201, 82, 233, 17, 152, 239, 179, 145, 124, 147, 33},
		{54, 216, 215, 148, 238, 156, 123, 242, 153, 51, 225, 50, 219, 99, 82, 194},
		{105, 71, 150, 88, 101, 21, 32, 136, 53, 28, 235, 133, 95, 114, 36, 110},
	}
)

func TestGoldenStructHashes(t *testing.T) {
	for i := range goldenStructs {
		t.Run(fmt.Sprintf("goldenStruct_%d", i), func(t *testing.T) {
			h, err := Hash(goldenStructs[i], FormatMD5, nil)
			if err != nil {
				t.Errorf("error hashing: %v", err)
			}
			if !bytes.Equal(h, goldenHashes[i]) {
				t.Errorf("incorrect hash %v", h)
			}
		})
	}
}
