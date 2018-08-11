package set

import (
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	b.ReportAllocs()

	s := New()
	for i := 0; i < b.N; i++ {
		s.Add("a")
	}
}

func TestAdd(t *testing.T) {
	s := New()

	s.Add("a")
	if size := s.Size(); size != 1 {
		t.Errorf("Expected a size of 1, got: %d", size)
	}
	if !s.Has("a") {
		t.Errorf("Expected \"a\" to be present")
	}

	s.Add("a")
	if size := s.Size(); size != 1 {
		t.Errorf("Expected a size of 1, got: %d", size)
	}

	s.Add("b")
	if size := s.Size(); size != 2 {
		t.Errorf("Expected a size of 2, got: %d", size)
	}
	if !s.Has("b") {
		t.Errorf("Expected \"b\" to be present")
	}
}

func BenchmarkRemove(b *testing.B) {
	b.ReportAllocs()

	s := New("a", "b", "c", "d")
	for i := 0; i < b.N; i++ {
		s.Remove("a")
	}
}

func TestRemove(t *testing.T) {
	s := New("a", "b", "c", "d")

	s.Remove("a")
	if size := s.Size(); size != 3 {
		t.Errorf("Expected a size of 3, got: %d", size)
	}
	if s.Has("a") {
		t.Errorf("Expected \"a\" to be absent")
	}

	s.Remove("a")
	if size := s.Size(); size != 3 {
		t.Errorf("Expected a size of 3, got: %d", size)
	}

	s.Remove("b")
	if size := s.Size(); size != 2 {
		t.Errorf("Expected a size of 2, got: %d", size)
	}
	if s.Has("b") {
		t.Errorf("Expected \"b\" to be absent")
	}

	s.Remove("e")
	if size := s.Size(); size != 2 {
		t.Errorf("Expected a size of 2, got: %d", size)
	}
}

func BenchmarkHas(b *testing.B) {
	b.ReportAllocs()

	s := New("a", "b", "c", "d")
	for i := 0; i < b.N; i++ {
		s.Has("a")
	}
}

func TestHas(t *testing.T) {
	s := New("a", "b", "c", "d")
	if !s.Has("a") {
		t.Errorf("Expected \"a\" to be present")
	}
	if !s.Has("b") {
		t.Errorf("Expected \"b\" to be present")
	}
	if s.Has("e") {
		t.Errorf("Expected \"e\" to be absent")
	}
}

func BenchmarkSize(b *testing.B) {
	b.ReportAllocs()

	s := New("a", "b", "c", "d")
	for i := 0; i < b.N; i++ {
		s.Size()
	}
}

func BenchmarkClear(b *testing.B) {
	b.ReportAllocs()

	s := New()
	for i := 0; i < b.N; i++ {
		s.Clear()
	}
}

func TestClear(t *testing.T) {
	s := New("a", "b", "c", "d")
	s.Clear()
	if size := s.Size(); size != 0 {
		t.Errorf("Expected a size of 0, got: %d", size)
	}
	if s.Has("a") {
		t.Errorf("Expected \"a\" to be absent")
	}
}

func BenchmarkStringSlice(b *testing.B) {
	b.ReportAllocs()

	s := New("a", "b", "c", "d")
	for i := 0; i < b.N; i++ {
		s.StringSlice()
	}
}

func TestStringSlice(t *testing.T) {
	s := New("a", "b", "c", "d")
	if size := s.Size(); size != 4 {
		t.Errorf("Expected a size of 4, got: %d", size)
	}
	if !s.Has("a") {
		t.Errorf("Expected \"a\" to be present")
	}

	s.Clear()
	if size := s.Size(); size != 0 {
		t.Errorf("Expected a size of 0, got: %d", size)
	}
	if s.Has("a") {
		t.Errorf("Expected \"a\" to be absent")
	}
}
