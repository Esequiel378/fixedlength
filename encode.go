package fixedlength

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrInvalidBooleanValue = errors.New("fixedlength: invalid boolean value")
	ErrInvalidIntValue     = errors.New("fixedlength: invalid int value")
	ErrInvalidUintValue    = errors.New("fixedlength: invalid uint value")
	ErrInvalidFloatValue   = errors.New("fixedlength: invalid float value")
	ErrUnsupportedKind     = errors.New("fixedlength: unsupported kind")
)

var (
	intBitSizes = map[reflect.Kind]int{
		reflect.Int:   0,
		reflect.Int8:  8,
		reflect.Int16: 16,
		reflect.Int32: 32,
		reflect.Int64: 64,
	}
	uintBitSizes = map[reflect.Kind]int{
		reflect.Uint:    0,
		reflect.Uint8:   8,
		reflect.Uint16:  16,
		reflect.Uint32:  32,
		reflect.Uint64:  64,
		reflect.Uintptr: 64,
	}
)

// setFieldValue sets the value for a struct field using reflection.
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := intBitSizes[field.Kind()]

		intVal, err := strconv.ParseInt(value, 10, bitSize)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bitSize := uintBitSizes[field.Kind()]

		uintVal, err := strconv.ParseUint(value, 10, bitSize)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidUintValue, err)
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		bitSize := func() int {
			if field.Kind() == reflect.Float32 {
				return 32
			}

			return 64
		}()

		floatVal, err := strconv.ParseFloat(value, bitSize)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidFloatValue, err)
		}
		field.SetFloat(floatVal)

	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidBooleanValue, err)
		}
		field.SetBool(boolVal)

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		err := setFieldValue(field.Elem(), value)
		if err != nil {
			return err
		}

		if field.Elem().IsZero() && field.CanSet() {
			field.Set(reflect.Zero(field.Type()))
		}
		return nil

	default:
		if field.CanAddr() && implementsUnmarshaler(field) {
			if um, ok := field.Addr().Interface().(Unmarshaler); ok {
				return um.Unmarshal([]byte(value))
			}
		}

		return fmt.Errorf("%w: %s", ErrUnsupportedKind, field.Kind())
	}

	return nil
}
