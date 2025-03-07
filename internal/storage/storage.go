package storage

import (
	"errors"
	"time"
)

var ErrEpmtyStorage = errors.New("empty storage")

type Storage interface {
	Push(interface{}, time.Time)
	GetElementsAt(time.Time) <-chan interface{}
	GetElements(int64) <-chan interface{}
	StoreAt() <-chan interface{}
	Show()
}
