package container

import (
	"sync"
)

type Queue struct {
	contents []interface{}
	cap int
	lock sync.Mutex
}

func NewQueue(cap int) (q *Queue) {
	defaultSize := cap*4
	if cap > 10240 {
		defaultSize = cap*2
	}

	return &Queue{
		contents: make([]interface{}, 0, defaultSize),
		cap:  cap,
	}
}

func (q *Queue) Size() int {
	return len(q.contents)
}

func (q *Queue) Enqueue(val interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.contents = append(q.contents, val)
}

func (q *Queue) ForceEnqueue(val interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.contents) > 0 {
		q.contents = q.contents[1:]
	}
	q.contents = append(q.contents, val)
}

func (q *Queue) Dequeue() (val interface{}) {
	if len(q.contents) < 1 {
		return nil
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	val = q.contents[0]
	return
}
