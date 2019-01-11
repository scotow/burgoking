package main

import (
	"fmt"
	"github.com/scotow/burgoking"
	"log"
)

func main() {
	gen := burgoking.NewCodeGenerator()

	code, err := gen.Generate()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(code)
}
