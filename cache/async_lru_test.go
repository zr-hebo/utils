package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_Async_LRU(t *testing.T) {
	cache := NewAsyncLRUCache(1, 10, 1*time.Second)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprint(i), fmt.Sprintf("test data %d", i))
		//time.Sleep(time.Millisecond * 50)
		a := 1000 - int(rand.Int31n(1000))

		val := cache.Get(fmt.Sprint(a))
		if val == nil {
			fmt.Printf("get %d --> nil\n", a)
		} else {
			fmt.Printf("get %d --> %v\n", a, val)
		}
	}
	time.Sleep(2 * time.Second)
}
