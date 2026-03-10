package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

const DefaultBufferSize int = 24 * 1024

func main() {

	var verbose bool
	var mtx sync.Mutex
	var wg sync.WaitGroup
	var workDir string = "."
	var total *big.Int = big.NewInt(0)

	if len(os.Args) > 1 {
		workDir = os.Args[1]

		if os.Args[1] == "--help" || os.Args[1] == "help" {
			fmt.Println("count lines in your project. blazingly fast\nUse VERBOSE=1 env variable for verbose output.")
			os.Exit(1)
		}
	}

	if os.Getenv("VERBOSE") == "1" {
		verbose = true
	} else {
		verbose = false
	}

	wg.Go(func() {
		filepath.WalkDir(workDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				wg.Go(func() {
					lines, err := countLines(path)
					if err != nil {
						log.Println("error: ", err.Error())
						return
					}

					if verbose {
						fmt.Println(filepath.Base(path), lines)
					}

					mtx.Lock()
					total = total.Add(total, big.NewInt(int64(lines)))
					mtx.Unlock()
				})
			}
			return nil
		})
	})

	wg.Wait()
	fmt.Println("total lines: ", total)
}

func countLines(path string) (int, error) {
	var lines int = 0
	var buf [DefaultBufferSize]byte
	var newline = []byte{'\n'}

	file, err := os.Open(path)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	for {
		n, err := file.Read(buf[:])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return lines, err
			}
		}

		lines += bytes.Count(buf[:n], newline)
	}

	return lines, nil
}
