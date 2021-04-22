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
        "key:%v, last visit time:%s", r.lruIndex.Value,
        r.lastVisitTime.Format("2006-01-02 15:04:05"))
}

// LRURecord record
type LRURecord struct {
    value         interface{}
    lruIndex      *list.Element
    lastVisitTime time.Time
}

// LRUCache lru cache
type LRUCache struct {
    // maxNum is the maximum number of cache entries before
    maxNum    int
    lock      sync.RWMutex
    orderList *list.List
    contents  map[string]*LRURecord
    TTL       int
}

// NewLRUCache create LRUCache instance
func NewLRUCache(num, ttl int) (rc *LRUCache) {
    if num < 1 {
        num = 0
    }

    return &LRUCache{
        maxNum:    num,
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
        record.value = val
        lc.orderList.MoveToFront(record.lruIndex)

    } else {
        // add new element
        if lc.orderList.Len() >= lc.maxNum {
            leastVisitIdx := lc.orderList.Back()
            delete(lc.contents, leastVisitIdx.Value.(string))
            lc.orderList.Remove(leastVisitIdx)
        }

        record = &LRURecord{
            value: val,
        }
        record.lruIndex = lc.orderList.PushFront(key)
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
            lc.orderList.Remove(record.lruIndex)
            delete(lc.contents, key)
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
func (lc *LRUCache) Remove(key string) (val interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    record, ok := lc.contents[key]
    if ok {
        lc.orderList.Remove(record.lruIndex)
        delete(lc.contents, key)
    }
    return
}

// Clean clean all value in cache
func (lc *LRUCache) Clean() (vals []interface{}) {
    oldContents := lc.contents
    oldOrderList := lc.orderList
    defer func() {
        vals = make([]interface{}, 0, len(oldContents))
        head := oldOrderList.Front()
        for head != nil {
            record := oldContents[head.Value.(string)]
            vals = append(vals, record.value)
            head = head.Next()
        }
    }()

    lc.lock.Lock()
    defer lc.lock.Unlock()
    lc.orderList = list.New()
    lc.contents = make(map[string]*LRURecord)
    return
}

// Clean clean last N value in cache
func (lc *LRUCache) CleanLast(num int) (vals []interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    totalNum := lc.orderList.Len()
    beginPos := totalNum - num
    if beginPos < 0 {
        beginPos = 0
    }

    vals = make([]interface{}, 0, beginPos)
    head := lc.orderList.Front()
    for head != nil {
        if beginPos > 0 {
            beginPos--
            head = head.Next()
            continue
        }

        recordKey := head.Value.(string)
        record := lc.contents[recordKey]
        vals = append(vals, record.value)

        delete(lc.contents, recordKey)
        oldHead := head
        head = head.Next()
        lc.orderList.Remove(oldHead)
    }
    return
}

// Clean clean values not change before given time in cache
func (lc *LRUCache) CleanBefore(sinceTime time.Time) (vals []interface{}) {
    lc.lock.Lock()
    defer lc.lock.Unlock()

    vals = make([]interface{}, 0, lc.orderList.Len())
    head := lc.orderList.Front()
    for head != nil {
        recordKey := head.Value.(string)
        record := lc.contents[recordKey]
        if record.lastVisitTime.After(sinceTime) {
            head = head.Next()
            continue
        }

        vals = append(vals, record.value)
        delete(lc.contents, recordKey)
        oldHead := head
        head = head.Next()
        lc.orderList.Remove(oldHead)
    }
    return
}

// Size size of cache
func (lc *LRUCache) Size() (size int) {
    lc.lock.RLock()
    size = len(lc.contents)
    lc.lock.RUnlock()
    return
}
