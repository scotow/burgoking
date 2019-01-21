package main

import (
	"flag"
	"fmt"
	"github.com/scotow/burgoking"
	"os"
	"sync"
)

var (
	count 		= flag.Int("n", 1, "number of code to generate")
	parallel 	= flag.Bool("p", true, "generate each code on a different goroutine")

	wg sync.WaitGroup
)

func generateCode() {
	code, err := burgoking.GenerateCode(nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(code)
	wg.Done()
}

func main() {
	flag.Parse()

	if *count <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid number of code")
		os.Exit(1)
	}

	if *parallel {
		for i := 0; i < *count; i++ {
			wg.Add(1)
			go generateCode()
		}

		wg.Wait()
	} else {
		for i := 0; i < *count; i++ {
			generateCode()
		}
	}
}
