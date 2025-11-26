package intmap

import (
	"maps"
	"math/rand"
	"slices"
	"testing"
	"unsafe"
)

func TestSet(t *testing.T) {
	s := NewSet[int](0)
	if sz := s.Len(); sz != 0 {
		t.Fatalf("length of set must be 0: %d", sz)
	}

	if added := s.Add(1); !added {
		t.Fatalf("1 must be added")
	}
	if sz := s.Len(); sz != 1 {
		t.Fatalf("length of set must be 1: %d", sz)
	}

	if added := s.Add(1); added {
		t.Fatalf("1 must not be added")
	}
	if sz := s.Len(); sz != 1 {
		t.Fatalf("length of set must be 1: %d", sz)
	}

	if added := s.Add(2); !added {
		t.Fatalf("2 must be added")
	}
	if sz := s.Len(); sz != 2 {
		t.Fatalf("length of set must be 2: %d", sz)
	}

	if !(s.Has(1) && s.Has(2)) {
		t.Fatalf("set must have both 1 and 2")
	}

	sum := 0
	s.ForEach(func(k int) bool {
		sum += k
		return true
	})
	if sum != 3 {
		t.Fatalf("total sum of elements must be 3")
	}
}

func TestSetClear(t *testing.T) {
	s := NewSet[int](0)

	s.Add(1)
	s.Add(2)
	if sz := s.Len(); sz != 2 {
		t.Fatalf("unexpected set len %d", sz)
	}
	if !s.Has(1) || !s.Has(2) {
		t.Fatalf("set must contain 1 and 2 before Clear()")
	}

	s.Clear()
	if sz := s.Len(); sz != 0 {
		t.Fatalf("unexpected set len %d", sz)
	}
	if s.Has(1) || s.Has(2) {
		t.Fatalf("set must not contain 1 or 2 after Clear()")
	}
}

func TestSetDel(t *testing.T) {
	s := NewSet[int](0)

	s.Add(1)
	s.Add(2)
	if sz := s.Len(); sz != 2 {
		t.Fatalf("unexpected set len %d", sz)
	}
	if !s.Has(1) || !s.Has(2) {
		t.Fatalf("set must contain 1 and 2 before Clear()")
	}

	if found := s.Del(27); found {
		t.Fatalf("set must not contain 27")
	}

	// Delete 2
	if found := s.Del(2); !found {
		t.Fatalf("set must contain 2 on delete")
	}
	if sz := s.Len(); sz != 1 {
		t.Fatalf("unexpected set len %d", sz)
	}
	if s.Has(2) {
		t.Fatalf("set must not contain 2 after Del(2)")
	}
	if found := s.Del(2); found {
		t.Fatalf("set must not contain 2 on seconb deletion")
	}

	// Delete 1
	if found := s.Del(1); !found {
		t.Fatalf("set must contain 1 on delete")
	}
	if sz := s.Len(); sz != 0 {
		t.Fatalf("unexpected set len %d", sz)
	}
	if s.Has(1) {
		t.Fatalf("set must not contain 1 after Del(1)")
	}
	if found := s.Del(1); found {
		t.Fatalf("set must not contain 1 on seconb deletion")
	}
}

func TestNilSet(t *testing.T) {
	var s *Set[int]

	if sz := s.Len(); sz != 0 {
		t.Fatalf("length of nil set must be 0: %d", sz)
	}

	if s.Has(0) || s.Has(1) {
		t.Fatalf("nil set must not have 0 or 1")
	}

	count := 0
	s.ForEach(func(k int) bool {
		count++
		return true
	})
	if count != 0 {
		t.Fatalf("total count of elements in nil set must be 0")
	}

}

func TestSetIter(t *testing.T) {
	s := NewSet[int](10)
	for i := 0; i < 100; i++ {
		s.Add(i)
	}

	sum := 0
	for k := range s.All() {
		sum += k
	}

	const sumTo99 = 99 * (99 + 1) / 2
	if sum != sumTo99 {
		t.Fatalf("unexpected sum: %d, want %d", sum, sumTo99)
	}
}

func TestSetMemory(t *testing.T) {
	if sz := unsafe.Sizeof(pair[int8, struct{}]{}); sz != unsafe.Sizeof(int8(0)) {
		t.Fatalf("unsafe size of pair[int32, struct{}] must be 1: %d", sz)
	}
	if sz := unsafe.Sizeof(pair[int16, struct{}]{}); sz != unsafe.Sizeof(int16(0)) {
		t.Fatalf("unsafe size of pair[int32, struct{}] must be 2: %d", sz)
	}
	if sz := unsafe.Sizeof(pair[int32, struct{}]{}); sz != unsafe.Sizeof(int32(0)) {
		t.Fatalf("unsafe size of pair[int32, struct{}] must be 4: %d", sz)
	}
	if sz := unsafe.Sizeof(pair[int64, struct{}]{}); sz != unsafe.Sizeof(int64(0)) {
		t.Fatalf("unsafe size of pair[int64, struct{}] must be 8: %d", sz)
	}
}

func BenchmarkSetAdd10K(b *testing.B) {
	rr := rand.New(rand.NewSource(12345))
	seen := make(map[int]struct{}, 10_000)
	for len(seen) < 10_000 {
		seen[int(rr.Int31())] = struct{}{}
	}
	keys := slices.AppendSeq(make([]int, 0, 10_000), maps.Keys(seen))
	b.ResetTimer()

	for range b.N {
		s := NewSet[int](10_000)
		for _, k := range keys {
			s.Add(k)
		}
	}
}

func BenchmarkSetHas10K(b *testing.B) {
	rr := rand.New(rand.NewSource(12345))
	hitKeys := make([]int, 10_000)
	s := NewSet[int](len(hitKeys))
	for i := range hitKeys {
		k := rr.Int()
		hitKeys[i] = k
		s.Add(k)
	}
	missKeys := make([]int, 0, 10_000)
	for len(missKeys) < len(hitKeys) {
		k := rr.Int()
		if !s.Has(k) {
			missKeys = append(missKeys, k)
		}
	}

	b.Run("all_miss", func(b *testing.B) {
		for range b.N {
			for _, k := range missKeys {
				_ = s.Has(k)
			}
		}
	})

	b.Run("all_hit", func(b *testing.B) {
		for range b.N {
			for _, k := range hitKeys {
				_ = s.Has(k)
			}
		}
	})

	b.Run("mixed_hit", func(b *testing.B) {
		for range b.N {
			for i := range hitKeys {
				_ = s.Has(hitKeys[i])
				_ = s.Has(missKeys[i])
			}
		}
	})
}
