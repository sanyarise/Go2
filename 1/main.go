package main

import (
	"fmt"
	"sync"
)

const n = 100

func main() {
	var (
		ex int
		wg = sync.WaitGroup{}
	)
	wg.Add(n)
	for i := 0; i < n; i += 1 {
		go func() {
			ex += 1
			fmt.Println(ex)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Well Done!")


}
