package container

import (
	"bytes"
	"container/list"
	"encoding/json"
	"strconv"
	"sync"
)

type KVItem struct {
	key string
	val interface{}
}

type SortedMap struct {
	lock     sync.RWMutex
	keyMap   map[string]*list.Element
	dataList *list.List
}

func NewSortedMap() (sm *SortedMap) {
	return &SortedMap{
		keyMap:   make(map[string]*list.Element),
		dataList: list.New(),
	}
}

func (sm *SortedMap) Exist(key string) (ok bool) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	_, ok = sm.keyMap[key]
	return
}

func (sm *SortedMap) Size() int {
	sm.lock.RLock()
	sm.lock.RUnlock()
	return len(sm.keyMap)
}

func (sm *SortedMap) Get(key string) (val interface{}, ok bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	elem := sm.keyMap[key]
	if elem == nil {
		return
	}

	ok = true
	kvItem := elem.Value.(*KVItem)
	val = kvItem.val
	return
}

func (sm *SortedMap) Set(key string, val interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	var kvItem *KVItem
	elem := sm.keyMap[key]
	if elem != nil {
		kvItem = elem.Value.(*KVItem)
		kvItem.val = val
		return
	}

	kvItem = &KVItem{
		key: key,
		val: val,
	}
	elem = sm.dataList.PushBack(kvItem)
	sm.keyMap[key] = elem
	return
}

func (sm *SortedMap) Remove(key string) (val interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	elem := sm.keyMap[key]
	if elem == nil {
		return
	}

	delete(sm.keyMap, key)
	kvItem := sm.dataList.Remove(elem).(*KVItem)
	val = kvItem.val
	return
}

func (sm *SortedMap) Walk(visit func(key string, val interface{}) (breakFor bool, err error)) (err error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		kvItem := elem.Value.(*KVItem)
		var breakFor bool
		breakFor, err = visit(kvItem.key, kvItem.val)
		if err != nil {
			return
		}

		if breakFor {
			break
		}
	}
	return
}

func (sm *SortedMap) MarshalJSON() ([]byte, error) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	var buf = bytes.NewBuffer(make([]byte, 0, 32))
	buf.WriteByte('{')
	sep := []byte{}
	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		buf.Write(sep)
		kvItem := elem.Value.(*KVItem)
		buf.WriteString(strconv.Quote(kvItem.key))
		buf.WriteByte(':')
		valInBytes, err := json.Marshal(kvItem.val)
		if err != nil {
			return nil, err
		}
		buf.Write(valInBytes)
		sep = []byte{','}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
