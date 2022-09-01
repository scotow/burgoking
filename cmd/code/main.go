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
	parallel 	= flag.Bool("p", false, "generate each code on a different goroutine")
)

func generateCode() {
	code, err := burgoking.GenerateCodeStatic(nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(code)
}

func generateCodeWait(wg *sync.WaitGroup) {
	generateCode()
	wg.Done()
}

func main() {
	flag.Parse()

	if *count <= 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid number of code")
		os.Exit(1)
	}

	if *parallel {
		var wg sync.WaitGroup

		for i := 0; i < *count; i++ {
			wg.Add(1)
			go generateCodeWait(&wg)
		}

		wg.Wait()
	} else {
		for i := 0; i < *count; i++ {
			generateCode()
		}
	}
}
