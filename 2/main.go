package main

import (
	"fmt"
	"sync"
)

func main() {
	var mtx sync.Mutex
	mtx.Lock()
	someFunc(&mtx)
	mtx.Lock()
	someFunc(&mtx)
}
func someFunc(mtx2 *sync.Mutex) {
	fmt.Println("ok")
	defer mtx2.Unlock()
}
