package symlock

import "sync"

type SymLock struct {
	mutex   sync.Mutex
	entries map[string]*symLockEntry
}

func (s *SymLock) WithMutex(symbol string, action func()) {
	e := s.acquireEntry(symbol)
	defer s.releaseEntry(symbol, e)

	action()
}

func (s *SymLock) acquireEntry(symbol string) *symLockEntry {
	s.mutex.Lock()

	e, ok := s.entries[symbol]
	if !ok {
		if s.entries == nil {
			s.entries = make(map[string]*symLockEntry)
		}

		e = &symLockEntry{}
		s.entries[symbol] = e
	}

	e.refCount++

	s.mutex.Unlock()

	e.Lock()
	return e
}

func (s *SymLock) releaseEntry(symbol string, e *symLockEntry) {
	e.Unlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	e.refCount--
	if e.refCount <= 0 {
		delete(s.entries, symbol)
	}
}
