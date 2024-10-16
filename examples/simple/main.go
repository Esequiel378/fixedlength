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
Noah Scott          19910601333222555999.99   
Amelia Ross         19861127666555444400.45   
`

type Person struct {
	FullName  string  `map:"0,20"`
	BirthDate string  `map:"20,28"`
	SSN       string  `map:"28,37"`
	Income    float64 `map:"37,-1"`
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
