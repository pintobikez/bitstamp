package robot

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteLog(w string) {

	// For more granular writes, open a file for writing.
	f, err := os.Create("results.log")
	fmt.Println("Error writing to log")

	defer f.Close()

	// A `WriteString` is also available.
	_, err := f.WriteString(w + "\n")
	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()
}
