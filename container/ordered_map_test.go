package container

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	sm := NewOrderedMap()
	vals := []int{1, 21, 27, 9, 15}
	for _, val := range vals {
		sm.Set(fmt.Sprintf("key_%v", val), val)
	}
	sm.Set("string_key", "MR'Wang")
	sm.Set("key with\"", "MR\"Wang")
	t.Log("set to OrderedMap OK")

	jsonVal, err := json.Marshal(sm)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(jsonVal))

	checkKey := "key_27"
	checkVal := 27
	val, ok := sm.Get(checkKey)
	if !ok {
		t.Fatalf("key %s not found", checkKey)
	}
	if !reflect.DeepEqual(val, checkVal) {
		t.Logf("get by key:'%s', expect %v, got %v", checkKey, checkVal, val)
	}
	t.Log("get from OrderedMap OK")

	val = sm.Remove(checkKey)
	if !reflect.DeepEqual(val, checkVal) {
		t.Logf("remove by key:'%s', expect %v, got %v", checkKey, checkVal, val)
	}

	_, ok = sm.Get(checkKey)
	if ok {
		t.Fatalf("key %s found after remove", checkKey)
	}
	t.Log("remove from OrderedMap OK")

	visitFunc := func(key string, val interface{}) (breakFor bool, err error) {
		t.Logf("visit key:%s, value:%v OK", key, val)
		return
	}
	err = sm.Walk(visitFunc)
	if err != nil {
		t.Fatalf("walk sorted map failed <-- %s", err.Error())
	}
	t.Log("visit items in OrderedMap OK")

	jsonVal, err = json.Marshal(sm)
	if err != nil {
		t.Fatal(err.Error())
	}

	receiver := make(map[string]interface{})
	err = json.Unmarshal(jsonVal, &receiver)
	if err != nil {
		t.Fatal(err.Error())
	}
	for key, val := range receiver {
		gotVal, _ := sm.Get(key)
		if fmt.Sprint(val) != fmt.Sprint(gotVal) {
			t.Fatalf("unmarshal to map get by key:'%s', expect %v, got %v", key, gotVal, val)
		}
	}
	t.Logf("%#v", receiver)
}
