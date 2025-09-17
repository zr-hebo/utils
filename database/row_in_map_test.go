package database

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/zr-hebo/utils/container"
)

func Test_getRecordInOrderedMapFromReceiver(t *testing.T) {
	type args struct {
		receiver []interface{}
		fields   []Field
	}

	receiveVal := sql.NullString{
		String: "1024",
		Valid:  true,
	}
	wantRecord := container.NewOrderedMap()
	wantRecord.Set("cnt", int64(1024))
	tests := []struct {
		name       string
		args       args
		wantRecord *container.OrderedMap
	}{
		{
			name: "common case",
			args: args{
				receiver: []interface{}{&receiveVal},
				fields: []Field{
					{
						Name: "cnt",
						Type: "int64",
					},
				},
			},
			wantRecord: wantRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRecord := getRecordInOrderedMapFromReceiver(tt.args.receiver, tt.args.fields); !compareOrderedMap(gotRecord, tt.wantRecord) {
				t.Errorf("getRecordInOrderedMapFromReceiver() = %v, want %v", gotRecord, tt.wantRecord)
			}
		})
	}
}

func compareOrderedMap(map1, map2 *container.OrderedMap) (equal bool) {
	equal = true
	_ = map1.Walk(func(key string, val interface{}) (breakFor bool, err error) {
		cmpVal, ok := map2.Get(key)
		if !ok {
			breakFor = true
			equal = false
			return
		}

		if !reflect.DeepEqual(val, cmpVal) {
			breakFor = true
			equal = false
		}
		return
	})
	return
}
