package cache

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func (r *AsyncLRURecord) String() string {
	return fmt.Sprintf(
		"key:%v, last visit time:%s", r.lruIndex.Value,
		r.lastVisitTime.Format("2006-01-02 15:04:05"))
}

// AsyncLRURecord record
type AsyncLRURecord struct {
	value         interface{}
	key           string
	lruIndex      *list.Element
	lastVisitTime time.Time
}

// Returns true if the item has expired.
func (r *AsyncLRURecord) Expired(ttl int, t *time.Time) bool {
	behindTime := time.Second * time.Duration(ttl)
	if t.After(r.lastVisitTime.Add(behindTime)) {
		return true
	}
	return false
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *AsyncLRUCache) {
	j.stop = make(chan bool)
	tick := time.Tick(j.Interval)
	for {
		select {
		case <-tick:
			c.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopJanitor(c *AsyncLRUCache) {
	c.janitor.stop <- true
}

func runJanitor(c *AsyncLRUCache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}
	c.janitor = j
	go j.Run(c)
}

// LRUCache lru cache
type AsyncLRUCache struct {
	// maxNum is the maximum number of cache entries before
	maxNum    int
	lock      sync.RWMutex
	orderList *list.List
	contents  map[string]*AsyncLRURecord
	TTL       int
	janitor   *janitor
}

// NewLRUCache create LRUCache instance
func NewAsyncLRUCache(num, ttl int, cleanupInterval time.Duration) (rc *AsyncLRUCache) {
	if num < 1 {
		num = 0
	}

	c := &AsyncLRUCache{
		maxNum:    num,
		orderList: list.New(),
		contents:  make(map[string]*AsyncLRURecord),
		TTL:       ttl,
	}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(c, stopJanitor)
	}
	return c
}

// Set add kv item to lru cache
func (lc *AsyncLRUCache) Set(key string, val interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		// update value in element
		record.value = val
		lc.orderList.MoveToFront(record.lruIndex)

	} else {
		//no need to delete last element in cache
		//add new element
		record = &AsyncLRURecord{
			value: val,
			key:   key,
		}
		record.lruIndex = lc.orderList.PushFront(record)
		lc.contents[key] = record
	}

	record.lastVisitTime = time.Now()
}

// Get get value by key
func (lc *AsyncLRUCache) Get(key string) (val interface{}) {
	lc.lock.RLock()
	defer lc.lock.RUnlock()

	record, ok := lc.contents[key]
	if ok {
		// check if is expired record
		now := time.Now()
		behindTime := time.Second * time.Duration(lc.TTL)
		if now.After(record.lastVisitTime.Add(behindTime)) {
			//no need delete from cache if expire
			return nil
		}

		// reorder key position
		lc.orderList.MoveToFront(record.lruIndex)
		val = record.value
		record.lastVisitTime = now
		return val
	}
	return nil
}

// Remove remove value by key
func (lc *AsyncLRUCache) Remove(key string) (val interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		lc.orderList.Remove(record.lruIndex)
		delete(lc.contents, key)
	}
	return
}

func (c *AsyncLRUCache) DeleteExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for k, v := range c.contents {
		if v.Expired(c.TTL, &now) {
			c.orderList.Remove(v.lruIndex)
			delete(c.contents, k)
		}
	}
	if c.orderList.Len() > c.maxNum {
		index := 0
		newOrderList := list.New()
		for i := c.orderList.Front(); i != nil; i = i.Next() {
			//remove left element
			//remove from content
			if index >= c.maxNum {
				delete(c.contents, i.Value.(*AsyncLRURecord).key)
			} else {
				//c.orderList.Remove(i)
				newOrderList.PushBack(i.Value.(*AsyncLRURecord))
			}
			index++
		}
		c.orderList = newOrderList
	}
}

// Size size of cache
func (lc *AsyncLRUCache) Size() (size int) {
	lc.lock.RLock()
	size = len(lc.contents)
	lc.lock.RUnlock()
	return
}
