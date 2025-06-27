package hashstructure

import (
	"bytes"
	"testing"

	"github.com/markphelps/optional"
)

type optionalStruct struct {
	Int    optional.Int
	String optional.String
	Bool   optional.Bool
}

type outerStruct struct {
	Int         int
	InnerStruct *optionalStruct
}

func TestOptional(t *testing.T) {
	s1 := &outerStruct{
		Int: 42,
		InnerStruct: &optionalStruct{
			Int:    optional.NewInt(3),
			String: optional.NewString("hello"),
		},
	}
	s2 := &outerStruct{
		Int: 42,
		InnerStruct: &optionalStruct{
			Int:  optional.NewInt(2),
			Bool: optional.NewBool(false),
		},
	}

	h1, err := Hash(s1, FormatMD5, nil)
	if err != nil {
		t.Errorf("error hashing s1: %v", err)
	}
	h2, err := Hash(s2, FormatMD5, nil)
	if err != nil {
		t.Errorf("error hashing s2: %v", err)
	}

	if bytes.Equal(h1, h2) {
		t.Error("hashes were equal and should have been different")
	}
}
