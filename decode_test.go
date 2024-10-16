package mapper

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type nested struct {
		A string `map:"0,1"`
		B int    `map:"2,3"`
		C float64
	}

	type testStruct struct {
		Nested nested `map:"0,4"`
		Empty  string
	}

	data := []byte("AB2D")
	var v testStruct
	err := Unmarshal(data, &v)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if v.Nested.A != "A" {
		t.Errorf("Expected v.Nested.A to be 'A', got '%s'", v.Nested.A)
	}

	if v.Nested.B != 2 {
		t.Errorf("Expected v.Nested.B to be 2, got %d", v.Nested.B)
	}

	if v.Nested.C != 0 {
		t.Errorf("Expected v.Nested.C to be 0, got %f", v.Nested.C)
	}

	if v.Empty != "" {
		t.Errorf("Expected v.Empty to be empty, got '%s'", v.Empty)
	}
}

func TestUnmarshalError(t *testing.T) {
	data := []byte("AB2D")

	t.Run("non-pointer int", func(t *testing.T) {
		var v int
		err := Unmarshal(data, v)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "map: Unmarshal(non-pointer int)" {
			t.Errorf("Expected error to be 'map: Unmarshal(non-pointer int)', got '%s'", err.Error())
		}
	})

	t.Run("nil *int", func(t *testing.T) {
		var v2 *int
		err := Unmarshal(data, v2)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "map: Unmarshal(nil *int)" {
			t.Errorf("Expected error to be 'map: Unmarshal(nil *int)', got '%s'", err.Error())
		}
	})

	t.Run("non-pointer struct {}", func(t *testing.T) {
		var v3 struct{}
		err := Unmarshal(data, v3)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "map: Unmarshal(non-pointer struct {})" {
			t.Errorf("Expected error to be 'map: Unmarshal(non-pointer struct {})', got '%s'", err.Error())
		}
	})

	t.Run("nil", func(t *testing.T) {
		err := Unmarshal(data, nil)
		if err == nil {
			t.Fatalf("Expected Unmarshal to fail")
		}

		if err.Error() != "map: Unmarshal(nil)" {
			t.Errorf("Expected error to be 'map: Unmarshal(non-pointer struct {})', got '%s'", err.Error())
		}
	})
}

// Define a type that implements the Unmarshaler interface
type CustomTime struct{}

var _ Unmarshaler = (*CustomTime)(nil)

func (ct *CustomTime) Unmarshal(data []byte) error {
	return nil
}

// Define a type that does NOT implement the Unmarshaler interface
type NotUnmarshaler struct{}

// Test suite for implementsUnmarshaler function
func TestImplementsUnmarshaler(t *testing.T) {
	// Case 1: Type that implements Unmarshaler
	t.Run("implements Unmarshaler", func(t *testing.T) {
		// Create a reflect.Value for a type that implements Unmarshaler
		ct := &CustomTime{}
		val := reflect.ValueOf(ct)

		if !implementsUnmarshaler(val) {
			t.Errorf("expected true, got false for CustomTime")
		}
	})

	// Case 2: Type that does NOT implement Unmarshaler
	t.Run("does not implement Unmarshaler", func(t *testing.T) {
		// Create a reflect.Value for a type that does not implement Unmarshaler
		notUnmarshaler := NotUnmarshaler{}
		val := reflect.ValueOf(notUnmarshaler)

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for NotUnmarshaler")
		}
	})

	// Case 3: Non-addressable value
	t.Run("non-addressable value", func(t *testing.T) {
		// Use a literal value, which is non-addressable
		val := reflect.ValueOf(42) // An int, which is non-addressable in this context

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for non-addressable int")
		}
	})

	// Case 4: Nil pointer that implements Unmarshaler
	t.Run("nil pointer implements Unmarshaler", func(t *testing.T) {
		// Use a nil pointer to CustomTime (which implements Unmarshaler)
		var ct *CustomTime
		val := reflect.ValueOf(ct)

		if !implementsUnmarshaler(val) {
			t.Errorf("expected true, got false for nil pointer to CustomTime")
		}
	})

	// Case 5: Nil pointer that does NOT implement Unmarshaler
	t.Run("nil pointer does not implement Unmarshaler", func(t *testing.T) {
		// Use a nil pointer to NotUnmarshaler (which does not implement Unmarshaler)
		var notUnmarshaler *NotUnmarshaler
		val := reflect.ValueOf(notUnmarshaler)

		if implementsUnmarshaler(val) {
			t.Errorf("expected false, got true for nil pointer to NotUnmarshaler")
		}
	})
}
