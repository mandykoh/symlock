package symlock

import (
	"hash/fnv"
	"runtime"
	"sync"
)

// SymLock is the interface for a symbolic lock.
//
// WithMutex executes the specified action using the given symbol for
// mutual exclusion. No two goroutines using the same SymLock and the same
// symbol will execute concurrently with each other.
type SymLock interface {
	WithMutex(symbol string, action func())
}

type symLock []sync.Mutex

func (s symLock) WithMutex(symbol string, action func()) {
	p := s.partitionForSymbol(symbol)
	p.Lock()
	defer p.Unlock()

	action()
}

func (s symLock) partitionForSymbol(symbol string) *sync.Mutex {
	hash := fnv.New32a()
	hash.Write([]byte(symbol))

	index := hash.Sum32() % uint32(len(s))
	return &s[index]
}

// New creates and returns a new SymLock with a default number of partitions,
// equal to the number of processors.
func New() SymLock {
	return NewWithPartitions(runtime.NumCPU())
}

// NewWithPartitions creates and returns a new SymLock with the specified
// number of partitions. The number of partitions effectively places an upper
// limit on the degree of concurrency.
func NewWithPartitions(n int) SymLock {
	s := make(symLock, n)
	return &s
}
