package main

import (
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/cmd"
	"os"
)


func main() {
	var err error

	for range Only.Once {
		err = cmd.Execute()
		if err != nil {
			break
		}

	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}
