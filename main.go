package main

import (
	"fmt"
	"os"
)

const DefaultBufferSize int = 24 * 1024

var verbose bool

func main() {

	var workDir string = "."

	if len(os.Args) > 1 {

		switch os.Args[1] {
		case "--help", "help":
			fmt.Println("count lines & size in your project. blazingly fast\nUse VERBOSE=1 env variable for verbose output.\n\t-By Aditya Gupta")
			os.Exit(1)
		case "--size", "size":
			runCountSize(workDir)
		default:
			workDir = os.Args[1]
		}
	}

	if os.Getenv("VERBOSE") == "1" {
		verbose = true
	} else {
		verbose = false
	}

	runCountLines(workDir)
}
