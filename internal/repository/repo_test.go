package repository

import (
	"testing"
)

func BenchmarkAtomicPut1(b *testing.B) {
	repo := MemoryRepo{}
	for i := 0; i < b.N; i++ {
		repo.Put("")
	}
}

func BenchmarkMutexPut1(b *testing.B) {
	repo := MemoryRepoMutex{}
	for i := 0; i < b.N; i++ {
		repo.Put("")
	}
}

func BenchmarkAtomicGet1(b *testing.B) {
	repo := MemoryRepo{}
	repo.Put("")
	for i := 0; i < b.N; i++ {
		repo.Get(1)
	}
}

func BenchmarkMutexGet1(b *testing.B) {
	repo := MemoryRepoMutex{}
	repo.Put("")
	for i := 0; i < b.N; i++ {
		repo.Get(0)
	}
}

func BenchmarkAtomicPutN(b *testing.B) {
	repo := MemoryRepo{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repo.Put("")
		}
	})
}

func BenchmarkMutexPutN(b *testing.B) {
	repo := MemoryRepoMutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			repo.Put("")
		}
	})
}

func BenchmarkAtomicGetN(b *testing.B) {
	repo := MemoryRepo{}
	b.RunParallel(func(pb *testing.PB) {
		repo.Put("")
		for pb.Next() {
			repo.Get(1)
		}
	})
}

func BenchmarkMutexGetN(b *testing.B) {
	repo := MemoryRepoMutex{}
	b.RunParallel(func(pb *testing.PB) {
		repo.Put("")
		for pb.Next() {
			repo.Get(0)
		}
	})
}

func BenchmarkAtomicMixedN(b *testing.B) {
	repo := MemoryRepo{}
	b.RunParallel(func(pb *testing.PB) {
		id := 0
		repo.Put("")
		for pb.Next() {
			if id%2 == 0 {
				repo.Get(1)
			} else {
				repo.Put("")
			}
			id += 1
		}
	})
}

func BenchmarkMutexMixedN(b *testing.B) {
	repo := MemoryRepoMutex{}
	b.RunParallel(func(pb *testing.PB) {
		id := 0
		repo.Put("")
		for pb.Next() {
			if id%2 == 0 {
				repo.Get(1)
			} else {
				repo.Put("")
			}
			id += 1
		}
	})
}
