package cache

import (
	"sync"
)

type BlockBuffer struct {
	bufLock      sync.RWMutex
	bufContent   map[int]chan []byte
	maxBlockNum  int
	maxBlockSize int
}

func NewBlockBuffer(maxBlockNum, maxBlockSize int) *BlockBuffer {
	return &BlockBuffer{
		bufContent:   make(map[int]chan []byte),
		maxBlockSize: maxBlockSize,
		maxBlockNum:  maxBlockNum,
	}
}

func (bb *BlockBuffer) AcquireBuffer(size int) (buf []byte) {
	bb.bufLock.RLock()
	bufQueue, ok := bb.bufContent[size]
	bb.bufLock.RUnlock()
	if !ok {
		bufQueue = make(chan []byte, bb.maxBlockNum)
		bb.bufLock.Lock()
		bb.bufContent[size] = bufQueue
		bb.bufLock.Unlock()
	}

	select {
	case buf = <-bufQueue:
	default:
		buf = make([]byte, 0, size)
	}
	return
}

func (bb *BlockBuffer) ReleaseBuffer(buf []byte) {
	buf = buf[:0]
	size := cap(buf)
	if size >= bb.maxBlockSize {
		// 不回收太大的缓存块，防止 OOM
		return
	}

	bb.bufLock.RLock()
	bufQueue, ok := bb.bufContent[size]
	bb.bufLock.RUnlock()
	if !ok {
		return
	}

	select {
	case bufQueue <- buf:
	default:
	}
}
