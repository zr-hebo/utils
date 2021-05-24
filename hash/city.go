package hash

import (
	"github.com/hungrybirder/cityhash"
	"github.com/siddontang/go/hack"
)

func CityHash(val []byte) uint64 {
	return cityhash.CityHash64(val, uint32(len(val)))
}

func CityHashFromString(val string) uint64 {
	return CityHash(hack.Slice(val))
}

func CityHashWithSeed(val []byte, seed uint64) uint64 {
	return cityhash.CityHash64WithSeed(val, uint32(len(val)), seed)
}

func CityHashWithSeedFromString(val string, seed uint64) uint64 {
	return cityhash.CityHash64WithSeed(hack.Slice(val), uint32(len(val)), seed)
}
