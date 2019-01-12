package main

import (
	"fmt"
	"github.com/scotow/burgoking"
	"os"
)

func main() {
	code, err := burgoking.GenerateCode(nil)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Println(code)
}
