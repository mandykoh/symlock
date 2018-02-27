package symlock

import (
	"sync"
	"testing"
)

func TestSymLock(t *testing.T) {

	t.Run(".WithMutex()", func(t *testing.T) {

		t.Run("synchronises on a symbol", func(t *testing.T) {
			var s = New()

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
				t.Errorf("Possible race condition detected; expected count of %d but got %d", 10000*workerCount, count)
			}
		})
	})
}
