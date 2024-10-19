package fixedlength

import (
	"errors"
	"reflect"
	"strings"
)

// Unmarshaler is the interface implemented by types
// that can unmarshal themselves.
// Unmarshal must copy the input data if it wishes
// to retain the data after returning.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

// implementsUnmarshaler checks if a field implements the Unmarshaler interface
func implementsUnmarshaler(val reflect.Value) bool {
	// If the value is invalid (e.g., a nil value), return false
	if !val.IsValid() {
		return false
	}

	// Check if the value itself implements json.Unmarshaler
	if val.Type().Implements(unmarshalerType) {
		return true
	}

	// If the value is addressable, check if the pointer to it implements json.Unmarshaler
	if val.CanAddr() {
		return val.Addr().Type().Implements(unmarshalerType)
	}

	// Otherwise, it does not implement the interface
	return false
}

// InvalidUnmarshalError describes an invalid argument passed to [Unmarshal].
// (The argument to [Unmarshal] must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "range: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "range: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "range: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parses the given string into the provided struct v.
// v must be a pointer to a struct, and its fields should be tagged with `range:"<start>,<end>"`
// where start and end are the lower and upper bounds of the segment in the string.
// Unmarshal will parse nested structs recursively.
func Unmarshal(data []byte, v any) error {
	// Validate that v is a pointer to a struct
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	// Iterate over struct fields to map segment names to fields
	for i := 0; i < rv.Elem().NumField(); i++ {
		field := rv.Elem().Field(i)

		// Recursively parse the struct
		if field.Kind() == reflect.Struct && !implementsUnmarshaler(field) {
			if err := Unmarshal(data, field.Addr().Interface()); err != nil {
				return err
			}

			continue
		}

		tag := rv.Elem().Type().Field(i).Tag.Get("range")

		start, end, err := parseTag(tag, len(data))
		if errors.Is(err, ErrTagEmpty) {
			continue
		}
		if err != nil {
			return err
		}

		value := strings.TrimSpace(string(data[start:end]))

		if err := setFieldValue(field, value); err != nil {
			return err
		}
	}

	return nil
}
