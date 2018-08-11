package set

import (
	"testing"
)

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		New("a", "b", "c", "d")
	}
}

func TestNew(t *testing.T) {
	s := New("a", "b", "c", "d")
	if size := s.Size(); size != 4 {
		t.Errorf("Expected a size of 4, got: %d", size)
	}
}
