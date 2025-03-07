package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	t.Run("change size storage", func(t *testing.T) {
		ms := New()
		sizeStart := ms.size
		ms.SetSize("self", ms.size+100)

		require.NotEqual(t, sizeStart, ms.size)
	})
	t.Run("date at", func(t *testing.T) {
		ms := New()
		data := struct{ some string }{some: "some"}
		for i := 0; i < 200; i++ {
			ms.Push(data, time.Now())
		}

		actC := 0
		for range ms.GetElementsAt(time.Now().Add(1 * time.Microsecond)) {
			actC++
		}
		require.Equal(t, 0, actC)
	})
	t.Run("storage size limit", func(t *testing.T) {
		dStart := time.Now()
		tSize := 50

		ms := New()
		ms.SetSize("self", int64(tSize))

		data := struct{ some string }{some: "some"}
		for i := 0; i < 200; i++ {
			ms.Push(data, time.Now())
		}

		actC := 0
		for range ms.GetElementsAt(dStart) {
			actC++
		}

		require.Equal(t, tSize, actC)
	})

	t.Run("storage parallel", func(t *testing.T) {
		t.Parallel()
		dStart := time.Now()
		tSize := 500

		ms1 := New()
		ms1.SetSize("self", int64(tSize))

		ms2 := New()
		ms2.SetSize("self", int64(tSize))

		data := struct{ some string }{some: "some"}

		wg := &sync.WaitGroup{}
		wg.Add(10)

		for w := 0; w < 10; w++ {
			w := w
			go func() {
				defer wg.Done()
				for i := 0; i < 50; i++ {
					if w == 0 {
						ms1.Push(data, time.Now())
					} else {
						ms2.Push(data, time.Now())
					}
				}
			}()
		}

		wg.Wait()

		actC1 := 0
		for range ms1.GetElementsAt(dStart) {
			actC1++
		}

		actC2 := 0
		for range ms2.GetElementsAt(dStart) {
			actC2++
		}
		require.Equal(t, 50, actC1)
		require.Equal(t, 450, actC2)
		require.Less(t, actC1, actC2)
	})
}
