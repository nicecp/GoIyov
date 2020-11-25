package singleflight

import "sync"

type caller struct {
	wg sync.WaitGroup
	response interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	buffer map[string]*caller
}

func (group *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	group.mu.Lock()
	if group.buffer == nil {
		group.buffer = make(map[string]*caller)
	}

	if caller, ok := group.buffer[key]; ok {
		group.mu.Unlock()
		caller.wg.Wait()
		return caller.response, caller.err
	}

	caller := new(caller)
	group.buffer[key] = caller
	group.mu.Unlock()

	caller.wg.Add(1)
	caller.response, caller.err = fn()
	caller.wg.Done()

	group.mu.Lock()
	delete(group.buffer, key)
	group.mu.Unlock()
	return caller.response, caller.err
}
