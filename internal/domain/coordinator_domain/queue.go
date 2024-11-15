package coordinator_domain

import (
	"container/list"
	"sync"
)

type SimpleQueue struct {
	lock  *sync.Mutex
	queue *list.List
}

func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{
		lock:  &sync.Mutex{},
		queue: list.New(),
	}
}

func (q *SimpleQueue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.queue.Len() > 1 {
		panic("queue length should be less than 1")
	}
	q.queue.PushBack(v)
}

func (q *SimpleQueue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.queue.Len() == 0 {
		return nil
	}
	e := q.queue.Front()
	q.queue.Remove(e)
	return e.Value
}

func (q *SimpleQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue.Len()
}
