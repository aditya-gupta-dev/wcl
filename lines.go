package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

func runCountLines(args *ArgsModel) {

	var mtx sync.Mutex
	var wg sync.WaitGroup
	var total *big.Int = big.NewInt(0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		filepath.WalkDir(args.WorkDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
				return nil
			}

			if !d.IsDir() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					lines, err := countLines(path)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error counting lines in %s: %v\n", path, err)
						return
					}

					if args.Verbose {
						fmt.Println(filepath.Base(path), lines)
					}

					mtx.Lock()
					total = total.Add(total, big.NewInt(int64(lines)))
					mtx.Unlock()
				}()
			}
			return nil
		})
	}()

	wg.Wait()
	fmt.Println(total)

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
