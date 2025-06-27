package hashstructure

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/markphelps/optional"
)

func (w *walker) visitOptional(v reflect.Value, typeName string) error {
	switch typeName {
	case "String":
		os := v.Interface().(optional.String)
		str := "string" + os.OrElse("") // use a prefix to distinguish any possible value from nil case
		if !os.Present() && !w.opts.ZeroNil {
			str = "nil"
		}
		_, err := fmt.Fprint(w.h, str)
		return err

	case "Error":
		os := v.Interface().(optional.Error)
		str := "error" + os.OrElse(errors.New("")).Error() // use a prefix to distinguish any possible value from nil case
		if !os.Present() && !w.opts.ZeroNil {
			str = "nil"
		}
		_, err := fmt.Fprint(w.h, str)
		return err

	case "Bool":
		ob := v.Interface().(optional.Bool)
		if w.opts.IgnoreZeroValue && !ob.OrElse(false) {
			return nil
		}
		str := "nil"
		if ob.Present() {
			str = fmt.Sprintf("%t", ob.OrElse(false))
		}
		if !ob.Present() && w.opts.ZeroNil {
			str = "false" // treat nil as false
		}
		_, err := fmt.Fprint(w.h, str)
		return err

	case "Int8":
		oi := v.Interface().(optional.Int8)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int8(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			// since no Go primitive numeric type is 3 bytes,
			// there can exist no int8, uint8, 16, etc. that has the
			// same bytes as the string "nil". Therefore, updating the hash state
			// by writing "nil" will be distinct from any binary.Write below,
			// which is what we want (distinguishing nil from 0 or any other "present" val)
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Byte":
		oi := v.Interface().(optional.Byte)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, byte(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Int16":
		oi := v.Interface().(optional.Int16)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int16(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Int32":
		oi := v.Interface().(optional.Int32)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int32(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Rune":
		oi := v.Interface().(optional.Rune)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, rune(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Int64":
		oi := v.Interface().(optional.Int64)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Int":
		oi := v.Interface().(optional.Int)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, int64(oi.OrElse(0)))

	case "Uint8":
		oi := v.Interface().(optional.Uint8)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, uint8(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Uint16":
		oi := v.Interface().(optional.Uint16)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, uint16(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Uint32":
		oi := v.Interface().(optional.Uint32)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, uint32(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Uint64":
		oi := v.Interface().(optional.Uint64)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, uint64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Uint":
		oi := v.Interface().(optional.Uint)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, uint64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, uint64(oi.OrElse(0)))

	case "Float32":
		oi := v.Interface().(optional.Float32)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, float32(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0.0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Float64":
		oi := v.Interface().(optional.Float64)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, float64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == 0.0 {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Complex64":
		oi := v.Interface().(optional.Complex64)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, complex64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == complex64(0) {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Complex128":
		oi := v.Interface().(optional.Complex128)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, complex128(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == complex128(0) {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, oi.OrElse(0))

	case "Uintptr":
		oi := v.Interface().(optional.Uintptr)
		if w.opts.ZeroNil && !oi.Present() {
			return binary.Write(w.h, binary.LittleEndian, int64(0))
		}
		if w.opts.IgnoreZeroValue && oi.OrElse(0) == uintptr(0) {
			return nil
		}
		if !oi.Present() {
			_, err := fmt.Fprint(w.h, "nil")
			return err
		}
		return binary.Write(w.h, binary.LittleEndian, int64(oi.OrElse(0)))
	}

	return fmt.Errorf("unsupported optional type: %s", typeName)
}
