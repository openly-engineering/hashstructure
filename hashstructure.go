package hashstructure

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"reflect"
	"time"
)

// HashOptions are options that are available for hashing.
type HashOptions struct {
	// TagName is the struct tag to look at when hashing the structure.
	// By default this is "hash".
	TagName string

	// ZeroNil is flag determining if nil pointer should be treated equal
	// to a zero value of pointed type. By default this is false.
	ZeroNil bool

	// IgnoreZeroValue is determining if zero value fields should be
	// ignored for hash calculation.
	IgnoreZeroValue bool

	// SlicesAsSets assumes that a `set` tag is always present for slices.
	// Default is false (in which case the tag is used instead)
	SlicesAsSets bool

	// UseStringer will attempt to use fmt.Stringer always. If the struct
	// doesn't implement fmt.Stringer, it'll fall back to trying usual tricks.
	// If this is true, and the "string" tag is also set, the tag takes
	// precedence (meaning that if the type doesn't implement fmt.Stringer, we
	// panic)
	UseStringer bool
}

// Format specifies the hashing process used. Different formats typically
// generate different hashes for the same value and have different properties.
type Format uint

const (
	// To disallow the zero value
	formatInvalid Format = iota

	// FormatV2 is the current recommended format and fixes the issues
	// noted in FormatV1.
	FormatV2

	formatMax // so we can easily find the end
)

// Hash returns the hash value of an arbitrary value.
//
// If opts is nil, then default options will be used. See HashOptions
// for the default values.
//
// The "format" is required and must be one of the format values defined
// by this library. You should probably just use "FormatV2". This allows
// generated hashes uses alternate logic to maintain compatibility with
// older versions.
//
// Notes on the value:
//
//   - Unexported fields on structs are ignored and do not affect the
//     hash value.
//
//   - Adding an exported field to a struct with the zero value will change
//     the hash value.
//
// For structs, the hashing can be controlled using tags. For example:
//
//	struct {
//	    Name string
//	    UUID string `hash:"ignore"`
//	}
//
// The available tag values are:
//
//   - "ignore" or "-" - The field will be ignored and not affect the hash code.
//
//   - "set" - The field will be treated as a set, where ordering doesn't
//     affect the hash code. This only works for slices.
//
//   - "string" - The field will be hashed as a string, only works when the
//     field implements fmt.Stringer
func Hash(v any, format Format, opts *HashOptions) (uint64, error) {
	// Validate our format
	if format <= formatInvalid || format >= formatMax {
		return 0, &ErrFormat{}
	}

	// Create default options
	if opts == nil {
		opts = &HashOptions{}
	}

	tagName := opts.TagName
	if tagName == "" {
		tagName = "hash"
	}

	// Create our walker and walk the structure
	w := &walker{
		format: format,
		h:      fnv.New64(),
		tag:    tagName,
		opts:   opts,
	}
	err := w.visit(reflect.ValueOf(v), nil)
	return 0 /*TODO*/, err
}

type walker struct {
	format Format
	h      hash.Hash64
	tag    string

	opts *HashOptions
}

type visitOpts struct {
	// Flags are a bitmask of flags to affect behavior of this visit
	Flags visitFlag

	// Information about the struct containing this field
	Struct      any
	StructField string
}

var timeType = reflect.TypeOf(time.Time{})

// visit visits a value recursively and updates w.h
func (w *walker) visit(v reflect.Value, opts *visitOpts) error {
	t := reflect.TypeOf(0)

	// Loop since these can be wrapped in multiple layers of pointers
	// and interfaces.
	for {
		// If we have an interface, dereference it. We have to do this up
		// here because it might be a nil in there and the check below must
		// catch that.
		if v.Kind() == reflect.Interface {
			v = v.Elem()
			continue
		}

		if v.Kind() == reflect.Ptr {
			if w.opts.ZeroNil {
				t = v.Type().Elem()
			}
			v = reflect.Indirect(v)
			continue
		}

		break
	}

	// If it is nil, treat it like a zero.
	if !v.IsValid() {
		v = reflect.Zero(t)
	}

	// Binary writing can use raw ints, we have to convert to
	// a sized-int, we'll choose the largest...
	switch v.Kind() {
	case reflect.Int:
		v = reflect.ValueOf(int64(v.Int()))
	case reflect.Uint:
		v = reflect.ValueOf(uint64(v.Uint()))
	case reflect.Bool:
		var tmp int8
		if v.Bool() {
			tmp = 1
		}
		v = reflect.ValueOf(tmp)
	}

	k := v.Kind()

	// We can shortcut numeric values by directly binary writing them
	if k >= reflect.Int && k <= reflect.Complex64 {
		// A direct hash calculation
		return binary.Write(w.h, binary.LittleEndian, v.Interface())
	}

	switch v.Type() {
	case timeType:
		b, err := v.Interface().(time.Time).MarshalBinary()
		if err != nil {
			return err
		}

		err = binary.Write(w.h, binary.LittleEndian, b)
		return err
	}

	switch k {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			err := w.visit(v.Index(i), nil)
			if err != nil {
				return err
			}

		}
		return nil

	case reflect.Map:

		var includeMap IncludableMap
		if opts != nil && opts.Struct != nil {
			if v, ok := opts.Struct.(IncludableMap); ok {
				includeMap = v
			}
		}

		// Build the hash for the map. We do this by XOR-ing all the key
		// and value hashes. This makes it deterministic despite ordering.

		for _, k := range v.MapKeys() {
			v := v.MapIndex(k)
			if includeMap != nil {
				incl, err := includeMap.HashIncludeMap(
					opts.StructField, k.Interface(), v.Interface())
				if err != nil {
					return err
				}
				if !incl {
					continue
				}
			}

			err := w.visit(k, nil)
			if err != nil {
				return err
			}
			err = w.visit(v, nil)
			if err != nil {
				return err
			}

		}

		return nil

	case reflect.Struct:
		parent := v.Interface()
		var include Includable
		if impl, ok := parent.(Includable); ok {
			include = impl
		}

		if impl, ok := parent.(Hashable); ok {
			h, err := impl.Hash()
			if err != nil {
				return err
			}
			_, err = w.h.Write([]byte(fmt.Sprintf("%d", h)))
			return err
		}

		// If we can address this value, check if the pointer value
		// implements our interfaces and use that if so.
		if v.CanAddr() {
			vptr := v.Addr()
			parentptr := vptr.Interface()
			if impl, ok := parentptr.(Includable); ok {
				include = impl
			}

			if impl, ok := parentptr.(Hashable); ok {
				h, err := impl.Hash()
				if err != nil {
					return err
				}
				_, err = w.h.Write([]byte(fmt.Sprintf("%d", h)))
				return err
			}
		}

		t := v.Type()
		err := w.visit(reflect.ValueOf(t.Name()), nil)
		if err != nil {
			return err
		}

		l := v.NumField()
		for i := 0; i < l; i++ {
			if innerV := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				var f visitFlag
				fieldType := t.Field(i)
				if fieldType.PkgPath != "" {
					// Unexported
					continue
				}

				tag := fieldType.Tag.Get(w.tag)
				if tag == "ignore" || tag == "-" {
					// Ignore this field
					continue
				}

				if w.opts.IgnoreZeroValue {
					if innerV.IsZero() {
						continue
					}
				}

				// if string is set, use the string value
				if tag == "string" || w.opts.UseStringer {
					if impl, ok := innerV.Interface().(fmt.Stringer); ok {
						innerV = reflect.ValueOf(impl.String())
					} else if tag == "string" {
						// We only show this error if the tag explicitly
						// requests a stringer.
						return &ErrNotStringer{
							Field: v.Type().Field(i).Name,
						}
					}
				}

				// Check if we implement includable and check it
				if include != nil {
					incl, err := include.HashInclude(fieldType.Name, innerV)
					if err != nil {
						return err
					}
					if !incl {
						continue
					}
				}

				switch tag {
				case "set":
					f |= visitFlagSet
				}

				err := w.visit(reflect.ValueOf(fieldType.Name), nil)
				if err != nil {
					return err
				}

				err = w.visit(innerV, &visitOpts{
					Flags:       f,
					Struct:      parent,
					StructField: fieldType.Name,
				})
				if err != nil {
					return err
				}

			}

		}

		return nil

	case reflect.Slice:
		/*
			// We have two behaviors here. If it isn't a set, then we just
			// visit all the elements. If it is a set, then we do a deterministic
			// hash code.
			var h uint64
			var set bool
			if opts != nil {
				set = (opts.Flags & visitFlagSet) != 0
			}
			l := v.Len()
			for i := 0; i < l; i++ {
				current, err := w.visit(v.Index(i), nil)
				if err != nil {
					return 0, err
				}

				if set || w.sets {
					h = hashUpdateUnordered(h, current)
				} else {
					h = hashUpdateOrdered(w.h, h, current)
				}
			}

			// Important: read the docs for hashFinishUnordered
			h = hashFinishUnordered(w.h, h)
		*/

		return errors.New("slice TODO")

	case reflect.String:
		// Directly hash
		_, err := w.h.Write([]byte(v.String()))
		return err

	default:
		return fmt.Errorf("unknown kind to hash: %s", k)
	}

}

// visitFlag is used as a bitmask for affecting visit behavior
type visitFlag uint

const (
	visitFlagInvalid visitFlag = iota
	visitFlagSet               = iota << 1
)
