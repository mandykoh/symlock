package symlock

import "sync"

type partition struct {
	mutex   sync.Mutex
	entries map[string]*symLockEntry
}
