package fixedlength

import (
	"errors"
	"reflect"
	"testing"
)

func TestSetFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		kind      reflect.Kind
		value     string
		wantValue interface{}
		wantErr   error
	}{
		// Int tests
		{
			name:      "valid int",
			kind:      reflect.Int,
			value:     "42",
			wantValue: int64(42),
			wantErr:   nil,
		},
		{
			name:      "invalid int",
			kind:      reflect.Int,
			value:     "invalid",
			wantValue: int64(0),
			wantErr:   ErrInvalidIntValue,
		},

		// Float tests
		{
			name:      "valid float",
			kind:      reflect.Float64,
			value:     "42.5",
			wantValue: float64(42.5),
			wantErr:   nil,
		},
		{
			name:      "invalid float",
			kind:      reflect.Float64,
			value:     "invalid",
			wantValue: float64(0),
			wantErr:   ErrInvalidFloatValue,
		},

		// String tests
		{
			name:      "valid string",
			kind:      reflect.String,
			value:     "hello",
			wantValue: "hello",
			wantErr:   nil,
		},

		// Bool tests
		{
			name:      "valid bool true",
			kind:      reflect.Bool,
			value:     "true",
			wantValue: true,
			wantErr:   nil,
		},
		{
			name:      "valid bool false",
			kind:      reflect.Bool,
			value:     "false",
			wantValue: false,
			wantErr:   nil,
		},
		{
			name:      "invalid bool",
			kind:      reflect.Bool,
			value:     "invalid",
			wantValue: false,
			wantErr:   ErrInvalidBooleanValue,
		},

		// Unsupported kind
		{
			name:      "unsupported kind",
			kind:      reflect.Slice,
			value:     "somevalue",
			wantValue: nil,
			wantErr:   ErrUnsupportedKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a zero value for the given type
			var field reflect.Value
			switch tt.kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field = reflect.New(reflect.TypeOf(int64(0))).Elem()
			case reflect.Float32, reflect.Float64:
				field = reflect.New(reflect.TypeOf(float64(0))).Elem()
			case reflect.String:
				field = reflect.New(reflect.TypeOf("")).Elem()
			case reflect.Bool:
				field = reflect.New(reflect.TypeOf(false)).Elem()
			default:
				field = reflect.New(reflect.TypeOf([]byte{})).Elem() // Unsupported kind
			}

			// Call setFieldValue and check for errors
			err := setFieldValue(field, tt.value)

			// Check for expected error
			if err != nil && tt.wantErr == nil {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error %v, got none", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Check the value set in the field (only for supported kinds)
			if err == nil {
				switch tt.kind {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if field.Int() != tt.wantValue {
						t.Errorf("expected int value %v, got %v", tt.wantValue, field.Int())
					}
				case reflect.Float32, reflect.Float64:
					if field.Float() != tt.wantValue {
						t.Errorf("expected float value %v, got %v", tt.wantValue, field.Float())
					}
				case reflect.String:
					if field.String() != tt.wantValue {
						t.Errorf("expected string value %v, got %v", tt.wantValue, field.String())
					}
				case reflect.Bool:
					if field.Bool() != tt.wantValue {
						t.Errorf("expected bool value %v, got %v", tt.wantValue, field.Bool())
					}
				}
			}
		})
	}
}
