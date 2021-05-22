package hash

import (
	"github.com/siddontang/go/hack"
	"hash/maphash"
)

var seed = maphash.MakeSeed()

func MapHash(val []byte) uint64 {
	var hash maphash.Hash
	hash.SetSeed(seed)
	hash.Write(val)
	return hash.Sum64()
}

func MapHashFromString(val string) uint64 {
	return MapHash(hack.Slice(val))
}
