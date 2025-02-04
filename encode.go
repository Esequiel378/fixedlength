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
	ErrInvalidFloatValue   = errors.New("fixedlength: invalid float value")
	ErrUnsupportedKind     = errors.New("fixedlength: unsupported kind")
)

// setFieldValue sets the value for a struct field using reflection.
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Int:
		intVal, err := strconv.ParseInt(value, 10, 10)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)
	case reflect.Int8:
		intVal, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)
	case reflect.Int16:
		intVal, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)
	case reflect.Int32:
		intVal, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)
	case reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetInt(intVal)
	case reflect.Uint:
		uintVal, err := strconv.ParseUint(value, 10, 10)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetUint(uintVal)
	case reflect.Uint8:
		uintVal, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetUint(uintVal)
	case reflect.Uint16:
		uintVal, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetUint(uintVal)
	case reflect.Uint32:
		uintVal, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetUint(uintVal)
	case reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return errors.Join(ErrInvalidIntValue, err)
		}
		field.SetUint(uintVal)
	case reflect.Float32:
		floatVal, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return errors.Join(ErrInvalidFloatValue, err)
		}
		field.SetFloat(floatVal)
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.Join(ErrInvalidFloatValue, err)
		}
		field.SetFloat(floatVal)
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return errors.Join(ErrInvalidBooleanValue, err)
		}
		field.SetBool(boolVal)
	case reflect.Ptr:
		// Create new value if pointer is nil
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		// Store original pointer state
		wasNil := field.IsNil()

		// Set the value on the dereferenced pointer
		err := setFieldValue(field.Elem(), value)
		if err != nil {
			return err
		}

		// Always nil the pointer if the element is zero (regardless of initial state)
		if field.Elem().IsZero() && field.CanSet() {
			field.Set(reflect.Zero(field.Type()))
		} else if wasNil && !field.Elem().IsZero() {
			// If we created a new value, ensure it's preserved
			field.Set(field)
		}

		return nil
	default:
		if implementsUnmarshaler(field) {
			um := field.Addr().Interface().(Unmarshaler)
			return um.Unmarshal([]byte(value))
		}

		return fmt.Errorf("%w: %s", ErrUnsupportedKind, field.Kind())
	}

	return nil
}
