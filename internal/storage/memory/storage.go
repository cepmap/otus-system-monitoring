package memorystorage

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
)

type element struct {
	timestamp time.Time
	data      interface{}
}

type MemoryStorage struct {
	rwm  sync.RWMutex
	list *list.List
	size int64
}

func New() *MemoryStorage {
	return &MemoryStorage{rwm: sync.RWMutex{}, list: list.New(), size: config.DaemonConfig.Stats.Limit + 1}
}

func (ms *MemoryStorage) SetSize(owner string, newsize int64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	ms.size = newsize
	logger.Info(fmt.Sprintf("[%s] changed size of storage. New size: %d", owner, newsize))
}

func (ms *MemoryStorage) Push(s interface{}, t time.Time) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	if ms.size == 0 {
		return
	}
	if ms.list.Len() == int(ms.size) {
		ms.list.Remove(ms.list.Back())
	}
	ms.list.PushFront(element{timestamp: t, data: s})
}

func (ms *MemoryStorage) GetElementsAt(t time.Time) <-chan interface{} {
	elemCh := make(chan interface{})
	go func() {
		ms.rwm.RLock()
		defer close(elemCh)
		defer ms.rwm.RUnlock()
		for last := ms.list.Front(); last != nil; last = last.Next() {
			elem := last.Value.(element)
			if t.After(elem.timestamp) {
				return
			}
			elemCh <- elem.data
		}
	}()

	return elemCh
}

func (ms *MemoryStorage) GetElements(num int64) <-chan interface{} {
	elemCh := make(chan interface{})
	go func() {
		ms.rwm.RLock()
		defer close(elemCh)
		defer ms.rwm.RUnlock()
		last := ms.list.Front()
		for num > 0 {
			if last == nil {
				break
			}
			elem := last.Value.(element)
			elemCh <- elem.data
			last = last.Next()
			num--
		}
	}()

	return elemCh
}

func (ms *MemoryStorage) Show() {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()

	for e := ms.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s: %+v\n", e.Value.(element).timestamp, e.Value.(element).data)
	}
}

func (ms *MemoryStorage) StoreAt() <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		ms.rwm.RLock()
		defer close(ch)
		defer ms.rwm.RUnlock()

		if ms.list.Len() == 0 {
			return
		}

		ll := ms.list.Back()
		if elem, ok := ll.Value.(element); ok {
			ch <- elem.timestamp
		}
	}()

	return ch
}

func (ms *MemoryStorage) GetTimestamp(value interface{}) (time.Time, bool) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()

	for e := ms.list.Front(); e != nil; e = e.Next() {
		elem := e.Value.(element)
		if elem.data == value {
			return elem.timestamp, true
		}
	}
	return time.Time{}, false
}

func (ms *MemoryStorage) Remove(value interface{}) bool {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	for e := ms.list.Front(); e != nil; e = e.Next() {
		elem := e.Value.(element)
		if elem.data == value {
			ms.list.Remove(e)
			return true
		}
	}
	return false
}
