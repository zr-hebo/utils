package hash

import (
	"github.com/siddontang/go/hack"
	"github.com/zentures/cityhash"
)

func CityHash(val []byte) uint64 {
	return cityhash.CityHash64WithSeed(val, 3, 0x57c5208e8f021a77)
}

func CityHashFromString(val string) uint64 {
	return CityHash(hack.Slice(val))
}

