package cache

import (
    "bytes"
    "container/list"
    "fmt"
    "sync"
    "time"
)

func (r *LRURecord) String() string {
    return fmt.Sprintf(
        "key:%v, last visit time:%s", r.key,
        r.lastVisitTime.Format("2006-01-02 15:04:05"))
}

// LRURecord record
type LRURecord struct {
    key  string
    node *list.Element
    // val           interface{}
    lastVisitTime time.Time
}

// LRUCache lru cache
type LRUCache struct {
    // maxNum is the maximum number of cache entries before
    maxNum           int
    lock             *sync.RWMutex
    orderList        *list.List
    contents map[string]*LRURecord
    TTL      int
}

// NewLRUCache create LRUCache instance
func NewLRUCache(num, ttl int) (rc *LRUCache) {
    if num < 1 {
        num = 0
    }

    return &LRUCache{
        maxNum:    num,
        lock:      &sync.RWMutex{},
        orderList: list.New(),
        contents:  make(map[string]*LRURecord),
        TTL:       ttl,
    }
}

func (lc LRUCache) String() string {
    var buf bytes.Buffer
    fmt.Fprintf(&buf, "there are %d record in LRU cache", lc.orderList.Len())
    for _, record := range lc.contents {
        fmt.Fprint(&buf, fmt.Sprintf("%s; ", record))
    }

    return buf.String()
}

// Set add kv item to lru cache
func (lc *LRUCache) Set(key string, val interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    record, ok := lc.contents[key]
    if ok {
        // update value in element
        record.node.Value = val
        lc.orderList.MoveToFront(record.node)

    } else {
        // add new element
        if lc.orderList.Len() >= lc.maxNum {
            leastUsedElement := lc.orderList.Back()
            // TODO
            delete(lc.contents, key)
            lc.orderList.Remove(leastUsedElement)
        }

        record = &LRURecord{}
        record.key = key
        record.node = lc.orderList.PushFront(key)
        record.node.Value = val
        lc.contents[key] = record
    }

    record.lastVisitTime = time.Now()
}

// Get get value by key
func (lc *LRUCache) Get(key string) (val interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    record, ok := lc.contents[key]
    if ok {
        // check if is expired record
        now := time.Now()
        behindTime := time.Second * time.Duration(lc.TTL)
        if now.After(record.lastVisitTime.Add(behindTime)) {
            lc.orderList.Remove(record.node)
            delete(lc.contents, key)
            return nil
        }

        // reorder key position
        lc.orderList.PushFront(key)
        val = record.node.Value
        record.lastVisitTime = now
        return val
    }
    return nil
}

// Remove remove value by key
func (lc *LRUCache) Remove(key string) (val interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    record, ok := lc.contents[key]
    if ok {
        lc.orderList.Remove(record.node)
        delete(lc.contents, key)
    }
    return
}
