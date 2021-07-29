package container

import (
	"bytes"
	"container/list"
	"encoding/json"
	"strconv"
	"sync"
)

type kvPairItem struct {
	key string
	val interface{}
}

type OrderedMap struct {
	lock     sync.RWMutex
	keyMap   map[string]*list.Element
	dataList *list.List
}

func NewOrderedMap() (sm *OrderedMap) {
	return &OrderedMap{
		keyMap:   make(map[string]*list.Element),
		dataList: list.New(),
	}
}

func (sm *OrderedMap) Exist(key string) (ok bool) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	_, ok = sm.keyMap[key]
	return
}

func (sm *OrderedMap) Size() int {
	sm.lock.RLock()
	sm.lock.RUnlock()
	return len(sm.keyMap)
}

func (sm *OrderedMap) Get(key string) (val interface{}, ok bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	elem := sm.keyMap[key]
	if elem == nil {
		return
	}

	ok = true
	kvItem := elem.Value.(*kvPairItem)
	val = kvItem.val
	return
}

func (sm *OrderedMap) Set(key string, val interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	var kvItem *kvPairItem
	elem := sm.keyMap[key]
	if elem != nil {
		kvItem = elem.Value.(*kvPairItem)
		kvItem.val = val
		return
	}

	kvItem = &kvPairItem{
		key: key,
		val: val,
	}
	elem = sm.dataList.PushBack(kvItem)
	sm.keyMap[key] = elem
	return
}

func (sm *OrderedMap) Remove(key string) (val interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	elem := sm.keyMap[key]
	if elem == nil {
		return
	}

	delete(sm.keyMap, key)
	kvItem := sm.dataList.Remove(elem).(*kvPairItem)
	val = kvItem.val
	return
}

func (sm *OrderedMap) Walk(visit func(key string, val interface{}) (breakFor bool, err error)) (err error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		kvItem := elem.Value.(*kvPairItem)
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

func (sm *OrderedMap) MarshalJSON() ([]byte, error) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	var buf = bytes.NewBuffer(make([]byte, 0, 32))
	buf.WriteByte('{')
	sep := []byte{}
	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		buf.Write(sep)
		kvItem := elem.Value.(*kvPairItem)
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
