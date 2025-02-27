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
Noah Scott          19910601333222555999.99   
Amelia Ross         19861127666555444400.45   
`

type Person struct {
	FullName  string  `range:"0,20"`
	BirthDate string  `range:"20,28"`
	SSN       string  `range:"28,37"`
	Income    float64 `range:"37,-1"`
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
