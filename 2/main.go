package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/profile"
)

const cnt = 100

func main() {
	defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	var (
		cntr int
		lock sync.Mutex
		wg   sync.WaitGroup
	)
	wg.Add(cnt)
	for i := 0; i < cnt; i += 1 {
		if i == 50 {
			//Усыпляю по одной миллисекунде 
			//чтобы лучше было видно запуск
			//планировщика
			time.Sleep(1*time.Millisecond)
			runtime.Gosched()
			time.Sleep(1*time.Millisecond)
		}
		go func() {
			defer wg.Done()
			
			lock.Lock()
			
			defer lock.Unlock()
			cntr += 1
		}()
	}
	wg.Wait()
	fmt.Println(cntr)
}
