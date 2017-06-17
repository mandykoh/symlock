package symlock

import (
	"strconv"
	"sync"
	"testing"
)

func TestSymLock(t *testing.T) {

	t.Run("WithMutex() synchronises on a symbol", func(t *testing.T) {
		var s SymLock

		var count int
		var condMutex = &sync.Mutex{}
		var startCond = sync.NewCond(condMutex)
		var readyGroup sync.WaitGroup
		var doneGroup sync.WaitGroup

		var workerCount = 32

		readyGroup.Add(workerCount)
		doneGroup.Add(workerCount)

		incrementer := func() {
			startCond.L.Lock()
			readyGroup.Done()
			startCond.Wait()
			startCond.L.Unlock()

			for i := 0; i < 10000; i++ {
				s.WithMutex("someSymbol", func() {
					count++
				})
			}

			doneGroup.Done()
		}

		for i := 0; i < workerCount; i++ {
			go incrementer()
		}

		readyGroup.Wait()
		startCond.Broadcast()
		doneGroup.Wait()

		if count != 10000*workerCount {
			t.Errorf("Possible race condition detected; expected count of %d but got %d", 1000*workerCount, count)
		}
	})

	t.Run("WithMutex() creates an entry for each symbol", func(t *testing.T) {
		var s SymLock

		var workerCount = 32

		var condMutex = &sync.Mutex{}
		var finishCond = sync.NewCond(condMutex)
		var readyGroup sync.WaitGroup
		var doneGroup sync.WaitGroup

		readyGroup.Add(workerCount)
		doneGroup.Add(workerCount)

		incrementer := func(n int) {
			s.WithMutex(strconv.Itoa(n), func() {
				finishCond.L.Lock()
				readyGroup.Done()
				finishCond.Wait()
				finishCond.L.Unlock()
			})
			doneGroup.Done()
		}

		for i := 0; i < workerCount; i++ {
			go incrementer(i + 1)
		}

		readyGroup.Wait()

		if count := len(s.entries); count != workerCount {
			t.Errorf("Expected %d entries to be active up but found %d", workerCount, count)
		}

		finishCond.Broadcast()
		doneGroup.Wait()
	})

	t.Run("WithMutex() cleans up old entries", func(t *testing.T) {
		var s SymLock

		var workerCount = 32

		var doneGroup sync.WaitGroup
		doneGroup.Add(workerCount)

		incrementer := func() {
			s.WithMutex("someSymbol", func() {})
			doneGroup.Done()
		}

		for i := 0; i < workerCount; i++ {
			go incrementer()
		}

		doneGroup.Wait()

		if count := len(s.entries); count != 0 {
			t.Errorf("Expected all entries to be cleaned up but found %d", count)
		}
	})
}
