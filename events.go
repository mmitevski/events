package events

import (
	"errors"
	"sync"
	"sync/atomic"
)

type EventFunc func(args ...any) error

var subscriptions map[string][]EventFunc = make(map[string][]EventFunc)
var initialized int32 = 0
var lock sync.RWMutex

func isInitialized() bool {
	return atomic.LoadInt32(&initialized) == 1
}

func Publish(name string, args ...any) error {
	// Use atomic to avoid locking
	if isInitialized() {
		if funcs, ok := subscriptions[name]; ok {
			for _, f := range funcs {
				if err := f(args...); err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return errors.New("the event subsystem is not initialized")
	}
}

func Subscribe(name string, eventFunc EventFunc) {
	lock.Lock()
	defer lock.Unlock()
	if isInitialized() {
		panic("event system is already initialized")
	}
	funcs, ok := subscriptions[name]
	if !ok {
		funcs = make([]EventFunc, 0)
	}
	subscriptions[name] = append(funcs, eventFunc)
}

func Initialize() {
	lock.Lock()
	defer lock.Unlock()
	atomic.StoreInt32(&initialized, 1)
}
