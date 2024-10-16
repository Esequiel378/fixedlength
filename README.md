# mapper

**mapper** is a Go library that unmarshals raw, unformatted text into Go structs by mapping data to specific ranges within a line. This is especially useful for parsing fixed-width fields in text files or other unstructured data formats where fields are defined by their byte positions rather than delimiters.

The library supports custom field unmarshalling via the `Unmarshaler` interface, allowing for flexible handling of nested structs or specific data types.

## Features

- Map raw text data into Go structs using struct tags to define the byte ranges.
- Supports nested structs and custom unmarshalling via the `Unmarshaler` interface.
- Recursive unmarshalling of embedded structs.
- Handles various types, including strings, integers, floats, and custom-defined types.

## Installation

You can install the library using Go modules:

```bash
go get -u github.com/esequiel378/mapper
```

## Getting Started

### Example 1: Basic Struct Unmarshalling

Consider an input file with raw data like this:

```
Olivia Parker       199703221112223331550.85   
Liam Evans          19891008444555666675.25   
Emma Ward           200307137778889991200.00  
```

To map this data into a Go struct, you define the struct with tags indicating the byte ranges for each field:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"mapper"
)

var input = `
Olivia Parker       199703221112223331550.85   
Liam Evans          19891008444555666675.25   
Emma Ward           200307137778889991200.00  
`

type Person struct {
	FullName  string  `map:"0,20"`
	BirthDate string  `map:"20,28"`
	SSN       string  `map:"28,37"`
	Income    float64 `map:"37,-1"`  // -1 means till the end of the line
}

func main() {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		var p Person
		err := mapper.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			log.Fatalf("Unmarshal failed: %v", err)
		}
		fmt.Printf("%+v\n", p)
	}
}
```

In this example, the struct tag `map:"start,end"` is used to indicate the byte range for each field.

### Example 2: Custom Unmarshaling with Nested Structs

If you need custom parsing logic, such as converting a date from `YYYYMMDD` format, you can implement the `Unmarshaler` interface.

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"mapper"
)

var input = `
Olivia Parker       199703221112223331550.85   
Liam Evans          19891008444555666675.25   
`

// Custom type for parsing birthdate
type PersonBirthDate struct {
	time.Time
}

var _ mapper.Unmarshaler = (*PersonBirthDate)(nil)

func (p *PersonBirthDate) Unmarshal(data []byte) error {
	birthDate, err := time.Parse("20060102", string(data))  // Parses date as YYYYMMDD
	if err != nil {
		return err
	}
	*p = PersonBirthDate{Time: birthDate}
	return nil
}

type Person struct {
	FullName  string          `map:"0,20"`
	BirthDate PersonBirthDate `map:"20,28"`
	SSN       string          `map:"28,37"`
	Income    float64         `map:"37,-1"`
}

func main() {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		var p Person
		err := mapper.Unmarshal(scanner.Bytes(), &p)
		if err != nil {
			log.Fatalf("Unmarshal failed: %v", err)
		}
		fmt.Printf("%+v\n", p)
	}
}
```

In this case, the `PersonBirthDate` struct implements the `Unmarshaler` interface to handle custom date parsing.

## Struct Tags

The struct field tags for `mapper` are in the format:

```go
`map:"start,end"`
```

- **start**: The starting byte index (inclusive).
- **end**: The ending byte index (exclusive). You can use `-1` to indicate that the field should take all remaining bytes from the `start` index until the end of the line.

## Custom Types and Unmarshaling

To handle more complex data types, you can implement the `Unmarshaler` interface for your custom types. The interface looks like this:

```go
type Unmarshaler interface {
    Unmarshal(data []byte) error
}
```

If a struct field implements this interface, `mapper` will call its `Unmarshal` method during the unmarshalling process, allowing you to define custom parsing logic for that field.

## Testing

You can run the tests for the `mapper` library with:

```bash
go test -v ./...
```

This will execute all tests in the project and give verbose output.

## Contributing

Contributions are welcome! Please submit a pull request with your improvements or open an issue to discuss potential changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.
