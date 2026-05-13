package main

import (
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

func runCountSize(workDir string) {
	var waitGroup sync.WaitGroup
	var mtx sync.Mutex
	var total *big.Int = big.NewInt(0)

	waitGroup.Go(func() {
		filepath.WalkDir(workDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				waitGroup.Go(func() {
					size, err := countSize(path)

					if err != nil {
						return
					}

					if verbose {
						fmt.Println(filepath.Base(path), size)
					}

					mtx.Lock()
					total = total.Add(total, big.NewInt(size))
					mtx.Unlock()
				})
			}
			return nil
		})
	})

	waitGroup.Wait()
	fmt.Println(formatBytes(total))
}

func countSize(absPath string) (int64, error) {
	file, err := os.Stat(absPath)
	if err != nil {
		return 0, err
	}

	return file.Size(), nil
}

func formatBytes(b *big.Int) string {

	f := new(big.Float).SetInt(b)

	kb := new(big.Float).SetInt(big.NewInt(1 << 10))
	mb := new(big.Float).SetInt(big.NewInt(1 << 20))
	gb := new(big.Float).SetInt(big.NewInt(1 << 30))
	tb := new(big.Float).SetInt(big.NewInt(1 << 40))

	if b.Cmp(new(big.Int).SetInt64(1<<40)) >= 0 {
		f.Quo(f, tb)
		return fmt.Sprintf("%.2f TB", f)
	} else if b.Cmp(new(big.Int).SetInt64(1<<30)) >= 0 {
		f.Quo(f, gb)
		return fmt.Sprintf("%.2f GB", f)
	} else if b.Cmp(new(big.Int).SetInt64(1<<20)) >= 0 {
		f.Quo(f, mb)
		return fmt.Sprintf("%.2f MB", f)
	} else if b.Cmp(new(big.Int).SetInt64(1<<10)) >= 0 {
		f.Quo(f, kb)
		return fmt.Sprintf("%.2f KB", f)
	}

	return fmt.Sprintf("%s Bytes", b.String())
}
