package genericSync

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
)

// mapInterface is the interface Map implements.
type mapInterface interface {
	Load(int) (int, bool)
	Store(key, value int)
	LoadOrStore(key, value int) (actual int, loaded bool)
	LoadAndDelete(key int) (value int, loaded bool)
	Delete(int)
	Swap(key, value int) (previous int, loaded bool)
	CompareAndSwap(key, old, new int) (swapped bool)
	CompareAndDelete(key, old int) (deleted bool)
	Range(func(key, value int) (shouldContinue bool))
}

type bench struct {
	setup func(*testing.B, mapInterface)
	perG  func(b *testing.B, pb *testing.PB, i int, m mapInterface)
}

func benchMap(b *testing.B, bench bench) {
	for _, m := range [...]mapInterface{&Map[int, int]{}} {
		b.Run(fmt.Sprintf("%T", m), func(b *testing.B) {
			m = reflect.New(reflect.TypeOf(m).Elem()).Interface().(mapInterface)
			if bench.setup != nil {
				bench.setup(b, m)
			}

			b.ResetTimer()

			var i int64
			b.RunParallel(func(pb *testing.PB) {
				id := int(atomic.AddInt64(&i, 1) - 1)
				bench.perG(b, pb, id*b.N, m)
			})
		})
	}

}

func BenchmarkLoadMostlyHits(b *testing.B) {
	const hits, misses = 1023, 1

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadOrStoreCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.LoadOrStore(0, 0)
			}
		},
	})
}
func BenchmarkLoadAndDeleteCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				if _, loaded := m.LoadAndDelete(0); loaded {
					m.Store(0, 0)
				}
			}
		},
	})
}

func BenchmarkRange(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Range(func(_, _ int) bool { return true })
			}
		},
	})
}

// BenchmarkAdversarialAlloc tests performance when we store a new value
// immediately whenever the map is promoted to clean and otherwise load a
// unique, missing key.
//
// This forces the Load calls to always acquire the map's mutex.
func BenchmarkAdversarialAlloc(b *testing.B) {
	benchMap(b, bench{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			var stores, loadsSinceStore int
			for ; pb.Next(); i++ {
				m.Load(i)
				if loadsSinceStore++; loadsSinceStore > stores {
					m.LoadOrStore(i, stores)
					loadsSinceStore = 0
					stores++
				}
			}
		},
	})
}

// BenchmarkAdversarialDelete tests performance when we periodically delete
// one key and add a different one in a large map.
//
// This forces the Load calls to always acquire the map's mutex and periodically
// makes a full copy of the map despite changing only one entry.
func BenchmarkAdversarialDelete(b *testing.B) {
	const mapSize = 1 << 10

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < mapSize; i++ {
				m.Store(i, i)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Load(i)

				if i%mapSize == 0 {
					m.Range(func(k, _ int) bool {
						m.Delete(k)
						return false
					})
					m.Store(i, i)
				}
			}
		},
	})
}

func BenchmarkDeleteCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Delete(0)
			}
		},
	})
}

func BenchmarkSwapCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Swap(0, 0)
			}
		},
	})
}

func BenchmarkSwapMostlyHits(b *testing.B) {
	const hits, misses = 1023, 1

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				if i%(hits+misses) < hits {
					v := i % (hits + misses)
					m.Swap(v, v)
				} else {
					m.Swap(i, i)
					m.Delete(i)
				}
			}
		},
	})
}

func BenchmarkSwapMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				if i%(hits+misses) < hits {
					v := i % (hits + misses)
					m.Swap(v, v)
				} else {
					m.Swap(i, i)
					m.Delete(i)
				}
			}
		},
	})
}

func BenchmarkCompareAndSwapCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for pb.Next() {
				if m.CompareAndSwap(0, 0, 42) {
					m.CompareAndSwap(0, 42, 0)
				}
			}
		},
	})
}

func BenchmarkCompareAndSwapNoExistingKey(b *testing.B) {
	benchMap(b, bench{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				if m.CompareAndSwap(i, 0, 0) {
					m.Delete(i)
				}
			}
		},
	})
}

func BenchmarkCompareAndSwapValueNotEqual(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.Store(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.CompareAndSwap(0, 1, 2)
			}
		},
	})
}

func BenchmarkCompareAndSwapMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				v := i
				if i%(hits+misses) < hits {
					v = i % (hits + misses)
				}
				m.CompareAndSwap(v, v, v)
			}
		},
	})
}

func BenchmarkCompareAndDeleteCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.LoadOrStore(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				if m.CompareAndDelete(0, 0) {
					m.Store(0, 0)
				}
			}
		},
	})
}

func BenchmarkCompareAndDeleteMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.LoadOrStore(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Load(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				v := i
				if i%(hits+misses) < hits {
					v = i % (hits + misses)
				}
				if m.CompareAndDelete(v, v) {
					m.Store(v, v)
				}
			}
		},
	})
}
