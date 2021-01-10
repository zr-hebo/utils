package cache

import (
    "bytes"
    "fmt"
    "sync"
    "time"
)

func (tr *TTLRecord) String() string {
    return fmt.Sprintf(
        "key:%v, last visit time:%s", tr.key,
        tr.lastVisitTime.Format("2006-01-02 15:04:05"))
}

// TTLRecord TTLRecord
type TTLRecord struct {
    key           string
    val           interface{}
    lastVisitTime time.Time
}

// TTLCache lru cache
type TTLCache struct {
    // maxNum is the maximum number of cache entries before
    maxNum   int
    lock     *sync.RWMutex
    contents map[string]*TTLRecord
    ttl      int
}

// NewLRUCache create TTLCache instance
func NewTTLCache(size, ttl int) (rc *TTLCache) {
    if size < 1 {
        size = 0
    }

    return &TTLCache{
        maxNum:   size,
        lock:     &sync.RWMutex{},
        contents: make(map[string]*TTLRecord),
        ttl:      ttl,
    }
}

func (tc TTLCache) String() string {
    var buf bytes.Buffer
    fmt.Fprintf(&buf, "there are %d TTLRecord in TTL cache", len(tc.contents))
    for _, tr := range tc.contents {
        fmt.Fprint(&buf, fmt.Sprintf("%s; ", tr))
    }

    return buf.String()
}

// Set add kv item to lru cache
func (tc *TTLCache) Set(key string, val interface{}) {
    tc.lock.Lock()
    defer tc.lock.Unlock()

    cr, ok := tc.contents[key]
    if ok {
        // update value in element
        cr.val = val

    } else {
        // add new element

        tc.contents[key] = cr
    }

    cr.lastVisitTime = time.Now()
}

// Get get value by key
func (tc *TTLCache) Get(key string) (val interface{}) {
    tc.lock.Lock()
    defer tc.lock.Unlock()

    CacheRecord, ok := tc.contents[key]
    if ok {
        // check if is expired TTLRecord
        now := time.Now()
        behindTime := time.Second * time.Duration(tc.ttl)
        if now.After(CacheRecord.lastVisitTime.Add(behindTime)) {
            delete(tc.contents, key)
            return nil
        }

        // reorder key position
        val = CacheRecord.val
        CacheRecord.lastVisitTime = now
        return val
    }
    return nil
}

// Remove remove value by key
func (tc *TTLCache) Remove(key string) (val interface{}) {
    tc.lock.Lock()
    defer tc.lock.Unlock()

    _, ok := tc.contents[key]
    if ok {
        delete(tc.contents, key)
    }
    return
}
