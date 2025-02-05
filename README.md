# fixedlength

**fixedlength** is a Go library that unmarshals fixed-width text into Go structs by mapping segments of a line (specified by byte ranges) into struct fields. This is especially useful when working with legacy systems, fixed-width data files, or any other unstructured formats where fields are defined by their positions rather than by delimiters.

The library supports custom field unmarshalling via the `Unmarshaler` interface, allowing for flexible handling of nested structs or special data types.

## Features

- **Field Mapping by Byte Ranges:** Use struct tags to map specific sections of a text line into corresponding struct fields.
- **Support for Basic and Custom Types:** Handles strings, integers, floats, booleans, and more. Custom parsing logic is possible by implementing the `Unmarshaler` interface.
- **Nested and Recursive Struct Parsing:** Automatically descends into nested structs if they do not implement custom unmarshalling.
- **Robust Error Handling:** Provides detailed errors (such as invalid numeric or boolean values and tag parsing errors) to facilitate easier debugging.
- **Flexible Tag Specification:** Use `-1` as the end index in a tag to capture all remaining bytes from the starting index.

## Struct Tags

The struct field tags for `fixedlength` follow this format:

```go
`range:"start,end"`
```

- **start**: The starting byte index (inclusive).
- **end**: The ending byte index (exclusive). Use -1 to indicate that the field should capture all remaining bytes from the start index until the end of the line.

Under the hood, these tags are parsed and validated. If the tag is empty or represents an inefectual (or reversed) range (e.g., start equals end or start is greater than end), a descriptive error is returned.

## Supported Data Types

The library supports the following data types out of the box:

- **String**: `string`
- **Integers**: `int`, `int8`, `int16`, `int32`, `int64`
- **Unsigned Integers**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`
- **Floating Point Numbers**: `float32`, `float64`
- **Boolean**: `bool`
- **Pointers**: Pointers to any of the above types are supported.
- **Custom Types**: Any type that implements the Unmarshaler interface can provide custom unmarshalling logic.

## Custom Types and Unmarshaling

For more advanced parsing, you can implement the Unmarshaler interface. This is especially useful for custom types such as dates, booleans, or any type that requires bespoke parsing logic.

```go
type Unmarshaler interface {
    Unmarshal(data []byte) error
}
```

When a field’s type implements this interface (or if its pointer does), fixedlength calls its Unmarshal method during the parsing process.

## Installation

Install the library using Go modules:

```bash
go get -u github.com/esequiel378/fixedlength
```

## Getting Started

### Example 1: Basic Struct Unmarshalling

Consider an input file with the following fixed-width records:

```
Olivia Parker       199703221112223331550.85
Liam Evans          19891008444555666675.25
Emma Ward           200307137778889991200.00
```

Define your struct with the corresponding range tags:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/esequiel378/fixedlength"
)

var input = `
Olivia Parker       199703221112223331550.85
Liam Evans          19891008444555666675.25
Emma Ward           200307137778889991200.00
`

type Person struct {
	FullName  string  `range:"0,20"`
	BirthDate string  `range:"20,28"`
	SSN       string  `range:"28,37"`
  Income    float64 `range:"37,-1"`  // -1 means read until the end of the line
}

func main() {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		var p Person
		err := fixedlength.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			log.Fatalf("Unmarshal failed: %v", err)
		}
		fmt.Printf("%+v\n", p)
	}
}
```

### Example 2: Custom Unmarshaling with Nested Structs

If you need to customize how a field is parsed (for example, converting a date from `YYYYMMDD`), implement the Unmarshaler interface:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/esequiel378/fixedlength"
)

var input = `
Olivia Parker       199703221112223331550.85
Liam Evans          19891008444555666675.25
`

// Custom type for parsing birthdate
type PersonBirthDate time.Time

var _ fixedlength.Unmarshaler = (*PersonBirthDate)(nil)

func (p *PersonBirthDate) Unmarshal(data []byte) error {
	birthDate, err := time.Parse("20060102", string(data))  // Parses date as YYYYMMDD
	if err != nil {
		return err
	}
	*p = PersonBirthDate(birthDate)
	return nil
}

type Person struct {
	FullName  string          `range:"0,20"`
	BirthDate PersonBirthDate `range:"20,28"`
	SSN       string          `range:"28,37"`
	Income    float64         `range:"37,-1"`
}

func main() {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		var p Person
		err := fixedlength.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			log.Fatalf("Unmarshal failed: %v", err)
		}
		fmt.Printf("%+v\n", p)
	}
}
```

## Testing

You can run the tests for the `fixedlength` library with:

```bash
go test -v ./...
```

This will execute all tests in the project and give verbose output.

## Contributing

Contributions are welcome! Please submit a pull request with your improvements or open an issue to discuss potential changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.
