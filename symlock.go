package symlock

import "hash/fnv"

type SymLock struct {
	partitions []partition
}

func (s *SymLock) WithMutex(symbol string, action func()) {
	e := s.acquireEntry(symbol)
	defer s.releaseEntry(symbol, e)

	action()
}

func (s *SymLock) acquireEntry(symbol string) *symLockEntry {
	p := s.partitionForSymbol(symbol)
	p.mutex.Lock()

	e, ok := p.entries[symbol]
	if !ok {
		if p.entries == nil {
			p.entries = make(map[string]*symLockEntry)
		}

		e = &symLockEntry{}
		p.entries[symbol] = e
	}

	e.refCount++

	p.mutex.Unlock()

	e.Lock()
	return e
}

func (s *SymLock) partitionForSymbol(symbol string) *partition {
	hash := fnv.New32a()
	hash.Write([]byte(symbol))

	index := hash.Sum32() % uint32(len(s.partitions))
	return &s.partitions[index]
}

func (s *SymLock) releaseEntry(symbol string, e *symLockEntry) {
	e.Unlock()

	p := s.partitionForSymbol(symbol)
	p.mutex.Lock()
	defer p.mutex.Unlock()

	e.refCount--
	if e.refCount <= 0 {
		delete(p.entries, symbol)
	}
}

func New() *SymLock {
	return NewWithPartitions(1)
}

func NewWithPartitions(n int) *SymLock {
	return &SymLock{
		partitions: make([]partition, n),
	}
}
