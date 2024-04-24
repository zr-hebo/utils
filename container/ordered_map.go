package container

import (
	"bytes"
	"container/list"
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
)

type KVPair struct {
	Key string
	Val interface{}
}

type OrderedMap struct {
	lock     sync.RWMutex
	keyMap   map[string]*list.Element
	dataList *list.List
}

func NewOrderedMap() (sm *OrderedMap) {
	return &OrderedMap{
		keyMap:   make(map[string]*list.Element, 16),
		dataList: list.New(),
	}
}

func NewOrderedMapWithSize(size int) (sm *OrderedMap) {
	return &OrderedMap{
		keyMap:   make(map[string]*list.Element, size),
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

func (sm *OrderedMap) Keys() []string {
	sm.lock.RLock()
	sm.lock.RUnlock()

	keys := make([]string, 0, len(sm.keyMap))
	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		kvItem := elem.Value.(*KVPair)
		keys = append(keys, kvItem.Key)
	}
	return keys
}

func (sm *OrderedMap) Pop() (kvItem *KVPair) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	elem := sm.dataList.Front()
	if elem != nil {
		sm.dataList.Remove(elem)
		kvItem = elem.Value.(*KVPair)
		delete(sm.keyMap, kvItem.Key)
	}
	return
}

func (sm *OrderedMap) Values() []interface{} {
	sm.lock.RLock()
	sm.lock.RUnlock()

	vals := make([]interface{}, 0, len(sm.keyMap))
	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		kvItem := elem.Value.(*KVPair)
		vals = append(vals, kvItem.Val)
	}
	return vals
}

func (sm *OrderedMap) Get(key string) (val interface{}, ok bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	elem := sm.keyMap[key]
	if elem == nil {
		return
	}

	ok = true
	kvItem := elem.Value.(*KVPair)
	val = kvItem.Val
	return
}

func (sm *OrderedMap) Set(key string, val interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	var kvItem *KVPair
	elem := sm.keyMap[key]
	if elem != nil {
		kvItem = elem.Value.(*KVPair)
		kvItem.Val = val
		return
	}

	kvItem = &KVPair{
		Key: key,
		Val: val,
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
	kvItem := sm.dataList.Remove(elem).(*KVPair)
	val = kvItem.Val
	return
}

func (sm *OrderedMap) Walk(visit func(key string, val interface{}) (breakFor bool, err error)) (err error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		kvItem := elem.Value.(*KVPair)
		var breakFor bool
		breakFor, err = visit(kvItem.Key, kvItem.Val)
		if err != nil {
			return
		}

		if breakFor {
			break
		}
	}
	return
}

func (sm *OrderedMap) String() string {
	if sm == nil {
		return "NULL"
	}

	val, err := sm.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(val)
}

func (sm *OrderedMap) MarshalJSON() ([]byte, error) {
	sm.lock.RLock()
	sm.lock.RUnlock()

	var buf = bytes.NewBuffer(make([]byte, 0, 32))
	buf.WriteByte('{')
	sep := []byte{}
	for elem := sm.dataList.Front(); elem != nil; elem = elem.Next() {
		buf.Write(sep)
		kvItem := elem.Value.(*KVPair)
		buf.WriteString(strconv.Quote(kvItem.Key))
		buf.WriteByte(':')
		valInBytes, err := json.Marshal(kvItem.Val)
		if err != nil {
			return nil, err
		}
		buf.Write(valInBytes)
		sep = []byte{','}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func DiffOrderedMap(from, to *OrderedMap) (fromDiff, toDiff *OrderedMap) {
	if to == nil {
		fromDiff = from
		return
	}
	if from == nil {
		toDiff = to
		return
	}

	fromDiff = NewOrderedMap()
	toDiff = NewOrderedMap()
	from.Walk(func(key string, fromVal interface{}) (breakFor bool, err error) {
		toVal, ok := to.Get(key)
		if !ok {
			fromDiff.Set(key, fromVal)
			return
		}

		if !reflect.DeepEqual(fromVal, toVal) {
			fromDiff.Set(key, fromVal)
			toDiff.Set(key, toVal)
		}
		return
	})

	to.Walk(func(key string, toVal interface{}) (breakFor bool, err error) {
		_, ok := from.Get(key)
		if !ok {
			toDiff.Set(key, toVal)
		}
		return
	})
	return
}
