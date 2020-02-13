package main

import "sync"

type group struct {
	sync.Mutex
	lastMsg     string
	repeatCount int
}

var (
	groups     = make(map[int64]*group)
	groupMutex sync.Mutex
)

func getGroup(id int64) *group {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	g, ok := groups[id]
	if !ok {
		g = new(group)
		groups[id] = g
	}
	return g
}
