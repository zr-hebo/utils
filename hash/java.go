package hash

import "github.com/siddontang/go/hack"

func JavaHash(val []byte) (hash int32) {
	if len(val) > 0 {
		for idx := 0; idx < len(val); idx++ {
			hash = 31*hash + int32(val[idx])
		}
	}
	return
}

func JavaHashFromString(val string) (hash int32) {
	return JavaHash(hack.Slice(val))
}

func Int32Abs(val int32) (abs int64) {
	if val >= 0 {
		abs = int64(val)
		return
	}
	return -1 * int64(val)
}

func Int64Abs(val int64) (abs int64) {
	if val >= 0 {
		abs = int64(val)
		return
	}
	return -1 * int64(val)
}
