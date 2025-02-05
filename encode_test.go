package fixedlength

import (
	"errors"
	"reflect"
	"testing"
)

var errError = errors.New("error")

type CustomUnmarshaler struct {
	Value string
}

var _ Unmarshaler = (*CustomUnmarshaler)(nil)

func (c *CustomUnmarshaler) Unmarshal(data []byte) error {
	if string(data) == "error" {
		return errError
	}

	c.Value = string(data)

	return nil
}

func createField(kind reflect.Kind, elemKind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.Ptr:
		return createPtrField(elemKind)
	case reflect.Slice:
		return reflect.ValueOf([]int{})
	default:
		return createBasicField(kind)
	}
}

func createBasicField(kind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.Uint:
		return reflect.New(reflect.TypeOf(uint(0))).Elem()
	case reflect.Uint8:
		return reflect.New(reflect.TypeOf(uint8(0))).Elem()
	case reflect.Uint16:
		return reflect.New(reflect.TypeOf(uint16(0))).Elem()
	case reflect.Uint32:
		return reflect.New(reflect.TypeOf(uint32(0))).Elem()
	case reflect.Uint64:
		return reflect.New(reflect.TypeOf(uint64(0))).Elem()
	case reflect.Int:
		return reflect.New(reflect.TypeOf(0)).Elem()
	case reflect.Int8:
		return reflect.New(reflect.TypeOf(int8(0))).Elem()
	case reflect.Int16:
		return reflect.New(reflect.TypeOf(int16(0))).Elem()
	case reflect.Int32:
		return reflect.New(reflect.TypeOf(int32(0))).Elem()
	case reflect.Int64:
		return reflect.New(reflect.TypeOf(int64(0))).Elem()
	case reflect.Float32:
		return reflect.New(reflect.TypeOf(float32(0))).Elem()
	case reflect.Float64:
		return reflect.New(reflect.TypeOf(float64(0))).Elem()
	case reflect.String:
		return reflect.New(reflect.TypeOf("")).Elem()
	case reflect.Bool:
		return reflect.New(reflect.TypeOf(false)).Elem()
	default:
		return reflect.Value{}
	}
}

func createPtrField(elemKind reflect.Kind) reflect.Value {
	switch elemKind {
	case reflect.Uint:
		return reflect.New(reflect.TypeOf(uint(0)))
	case reflect.Uint8:
		return reflect.New(reflect.TypeOf(uint8(0)))
	case reflect.Uint16:
		return reflect.New(reflect.TypeOf(uint16(0)))
	case reflect.Uint32:
		return reflect.New(reflect.TypeOf(uint32(0)))
	case reflect.Uint64:
		return reflect.New(reflect.TypeOf(uint64(0)))
	case reflect.Int:
		return reflect.New(reflect.TypeOf(0))
	case reflect.Int8:
		return reflect.New(reflect.TypeOf(int8(0)))
	case reflect.Int16:
		return reflect.New(reflect.TypeOf(int16(0)))
	case reflect.Int32:
		return reflect.New(reflect.TypeOf(int32(0)))
	case reflect.Int64:
		return reflect.New(reflect.TypeOf(int64(0)))
	case reflect.Float32:
		return reflect.New(reflect.TypeOf(float32(0)))
	case reflect.Float64:
		return reflect.New(reflect.TypeOf(float64(0)))
	case reflect.String:
		return reflect.New(reflect.TypeOf(""))
	case reflect.Bool:
		return reflect.New(reflect.TypeOf(false))
	case reflect.Struct:
		return reflect.New(reflect.TypeOf(CustomUnmarshaler{}))
	default:
		return reflect.ValueOf(nil)
	}
}

func TestSetFieldValue(t *testing.T) {
	tests := []struct {
		name          string
		kind          reflect.Kind
		elemKind      reflect.Kind
		value         string
		wantElemValue any
		wantErr       error
	}{
		// Signed Integer types
		{
			name:          "valid_int",
			kind:          reflect.Int,
			value:         "42",
			wantElemValue: 42,
		},
		{
			name:          "valid_int8",
			kind:          reflect.Int8,
			value:         "127",
			wantElemValue: int8(127),
		},
		{
			name:          "invalid_int8_overflow",
			kind:          reflect.Int8,
			value:         "128",
			wantElemValue: int8(0),
			wantErr:       ErrInvalidIntValue,
		},
		{
			name:          "valid_int16",
			kind:          reflect.Int16,
			value:         "32767",
			wantElemValue: int16(32767),
		},
		{
			name:          "invalid_int16_overflow",
			kind:          reflect.Int16,
			value:         "32768",
			wantElemValue: int16(0),
			wantErr:       ErrInvalidIntValue,
		},
		{
			name:          "valid_int32",
			kind:          reflect.Int32,
			value:         "2147483647",
			wantElemValue: int32(2147483647),
		},
		{
			name:          "invalid_int32_overflow",
			kind:          reflect.Int32,
			value:         "2147483648",
			wantElemValue: int32(0),
			wantErr:       ErrInvalidIntValue,
		},
		{
			name:          "valid_int64",
			kind:          reflect.Int64,
			value:         "42",
			wantElemValue: int64(42),
		},
		{
			name:          "invalid_int64",
			kind:          reflect.Int64,
			value:         "invalid",
			wantElemValue: int64(0),
			wantErr:       ErrInvalidIntValue,
		},

		// Unsigned Integer types
		{
			name:          "valid_uint",
			kind:          reflect.Uint,
			value:         "42",
			wantElemValue: uint(42),
		},
		{
			name:          "invalid_uint_negative",
			kind:          reflect.Uint,
			value:         "-42",
			wantElemValue: uint(0),
			wantErr:       ErrInvalidUintValue,
		},
		{
			name:          "valid_uint8",
			kind:          reflect.Uint8,
			value:         "255",
			wantElemValue: uint8(255),
		},
		{
			name:          "invalid_uint8_overflow",
			kind:          reflect.Uint8,
			value:         "256",
			wantElemValue: uint8(0),
			wantErr:       ErrInvalidUintValue,
		},
		{
			name:          "valid_uint16",
			kind:          reflect.Uint16,
			value:         "65535",
			wantElemValue: uint16(65535),
		},
		{
			name:          "invalid_uint16_overflow",
			kind:          reflect.Uint16,
			value:         "65536",
			wantElemValue: uint16(0),
			wantErr:       ErrInvalidUintValue,
		},
		{
			name:          "valid_uint32",
			kind:          reflect.Uint32,
			value:         "4294967295",
			wantElemValue: uint32(4294967295),
		},
		{
			name:          "invalid_uint32_overflow",
			kind:          reflect.Uint32,
			value:         "4294967296",
			wantElemValue: uint32(0),
			wantErr:       ErrInvalidUintValue,
		},
		{
			name:          "valid_uint64",
			kind:          reflect.Uint64,
			value:         "18446744073709551615",
			wantElemValue: uint64(18446744073709551615),
		},
		{
			name:          "invalid_uint64_overflow",
			kind:          reflect.Uint64,
			value:         "18446744073709551616",
			wantElemValue: uint64(0),
			wantErr:       ErrInvalidUintValue,
		},
		{
			name:          "valid_uintptr",
			kind:          reflect.Uintptr,
			value:         "12345",
			wantElemValue: uintptr(0),
			wantErr:       ErrUnsupportedKind,
		},
		{
			name:          "invalid_uintptr",
			kind:          reflect.Uintptr,
			value:         "invalid",
			wantElemValue: uintptr(0),
			wantErr:       ErrUnsupportedKind,
		},

		// Floating-point types
		{
			name:          "valid_float32",
			kind:          reflect.Float32,
			value:         "3.14",
			wantElemValue: float32(3.14),
		},
		{
			name:          "invalid_float32_non_numeric",
			kind:          reflect.Float32,
			value:         "abc",
			wantElemValue: float32(0),
			wantErr:       ErrInvalidFloatValue,
		},
		{
			name:          "float32_overflow",
			kind:          reflect.Float32,
			value:         "1e100",
			wantElemValue: float32(0),
			wantErr:       ErrInvalidFloatValue,
		},
		{
			name:          "valid_float64",
			kind:          reflect.Float64,
			value:         "42.5",
			wantElemValue: 42.5,
		},
		{
			name:          "invalid_float64",
			kind:          reflect.Float64,
			value:         "invalid",
			wantElemValue: 0.0,
			wantErr:       ErrInvalidFloatValue,
		},

		// String types
		{
			name:          "valid_string",
			kind:          reflect.String,
			value:         "hello",
			wantElemValue: "hello",
		},
		{
			name:          "string_empty",
			kind:          reflect.String,
			value:         "",
			wantElemValue: "",
		},
		{
			name:          "string_whitespace",
			kind:          reflect.String,
			value:         "   ",
			wantElemValue: "   ",
		},
		{
			name:          "string_special_chars",
			kind:          reflect.String,
			value:         "Hello\nWorld\t!",
			wantElemValue: "Hello\nWorld\t!",
		},

		// Boolean types
		{
			name:          "valid_bool_true",
			kind:          reflect.Bool,
			value:         "true",
			wantElemValue: true,
		},
		{
			name:          "valid_bool_false",
			kind:          reflect.Bool,
			value:         "false",
			wantElemValue: false,
		},
		{
			name:          "invalid_bool",
			kind:          reflect.Bool,
			value:         "invalid",
			wantElemValue: false,
			wantErr:       ErrInvalidBooleanValue,
		},
		{
			name:          "bool_1",
			kind:          reflect.Bool,
			value:         "1",
			wantElemValue: true,
		},
		{
			name:          "bool_0",
			kind:          reflect.Bool,
			value:         "0",
			wantElemValue: false,
		},

		// Pointer types
		{
			name:          "pointer_int_valid",
			kind:          reflect.Ptr,
			elemKind:      reflect.Int64,
			value:         "123",
			wantElemValue: int64(123),
		},
		{
			name:          "pointer_uint_valid",
			kind:          reflect.Ptr,
			elemKind:      reflect.Uint,
			value:         "42",
			wantElemValue: uint(42),
		},
		{
			name:          "pointer_float32_valid",
			kind:          reflect.Ptr,
			elemKind:      reflect.Float32,
			value:         "3.14",
			wantElemValue: float32(3.14),
		},
		{
			name:          "pointer_string_empty",
			kind:          reflect.Ptr,
			elemKind:      reflect.String,
			value:         "",
			wantElemValue: "",
		},
		{
			name:          "pointer_bool_valid_true",
			kind:          reflect.Ptr,
			elemKind:      reflect.Bool,
			value:         "true",
			wantElemValue: true,
		},

		// Custom Unmarshalers
		{
			name:          "custom_unmarshaler",
			kind:          reflect.Ptr,
			elemKind:      reflect.Struct,
			value:         "test",
			wantElemValue: CustomUnmarshaler{Value: "test"},
		},
		{
			name:          "custom_unmarshaler_error",
			kind:          reflect.Ptr,
			elemKind:      reflect.Struct,
			value:         "error",
			wantElemValue: CustomUnmarshaler{Value: ""},
			wantErr:       errError,
		},

		// Unsupported types
		{
			name:    "unsupported_complex64",
			kind:    reflect.Complex64,
			value:   "1+2i",
			wantErr: ErrUnsupportedKind,
		},
		{
			name:    "unsupported_chan",
			kind:    reflect.Chan,
			value:   "test",
			wantErr: ErrUnsupportedKind,
		},
		{
			name:    "unsupported_func",
			kind:    reflect.Func,
			value:   "test",
			wantErr: ErrUnsupportedKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := createField(tt.kind, tt.elemKind)
			err := setFieldValue(field, tt.value)

			if !errors.Is(err, tt.wantErr) && (err != nil || tt.wantErr != nil) {
				t.Errorf("error mismatch\ngot:  %v\nwant: %v", err, tt.wantErr)
			}

			if tt.kind == reflect.Ptr {
				if field.Kind() != reflect.Ptr {
					t.Errorf("expected pointer, got %v", field.Kind())
				}

				if !field.IsNil() {
					deref := field.Elem().Interface()
					if !reflect.DeepEqual(deref, tt.wantElemValue) {
						t.Errorf("value mismatch\ngot:  %v\nwant: %v", deref, tt.wantElemValue)
					}
				}
				return
			}

			if tt.wantErr == nil && tt.kind != reflect.Slice {
				got := field.Interface()
				if !reflect.DeepEqual(got, tt.wantElemValue) {
					t.Errorf("value mismatch\ngot:  %v\nwant: %v", got, tt.wantElemValue)
				}
			}
		})
	}
}
