package utils

import "sync/atomic"

type ClientManager struct {
	clientIDs map[atomic.Uint64]bool
	lastID    atomic.Uint64
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clientIDs: make(map[atomic.Uint64]bool),
	}
}

func (manager *ClientManager) NextClientID() {
	// newID := manager.lastID.Add(1)
	// return newID
}
