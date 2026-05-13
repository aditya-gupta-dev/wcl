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

func runCountLines(workDir string) {

	var mtx sync.Mutex
	var wg sync.WaitGroup
	var total *big.Int = big.NewInt(0)
	wg.Go(func() {
		filepath.WalkDir(workDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				wg.Go(func() {
					lines, err := countLines(path)
					if err != nil {
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
