package storage

import (
	"errors"
	"time"
)

var ErrEpmtyStorage = errors.New("empty storage")

type Storage interface {
	Push(item interface{}, timestamp time.Time)
	GetElementsAt(from time.Time) <-chan interface{}
	GetTimestamp(item interface{}) (time.Time, bool)
	Remove(item interface{}) bool
	GetElements(int64) <-chan interface{}
	StoreAt() <-chan interface{}
	Show()
}
