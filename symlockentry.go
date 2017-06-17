package symlock

import "sync"

type symLockEntry struct {
	sync.Mutex
	refCount uint32
}
