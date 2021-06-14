package collections

import (
	"sync"
)

type StringSet struct {
	mu sync.Mutex

	currIdx int
	m       map[int]string
	rm      map[string]int
	k       []string
}

func NewStringSet() *StringSet {
	return &StringSet{
		currIdx: 0,
		m:       make(map[int]string),
		rm:      make(map[string]int),
		k:       []string{},
	}
}

func (s *StringSet) Get(value string) (string, bool) {
	if s.Contains(value) {
		return value, true
	}

	return "N", false
}

func (s *StringSet) Add(value string) {
	if s.Contains(value) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[s.currIdx] = value
	s.rm[value] = s.currIdx
	s.k = append(s.k, value)
	s.currIdx += 1

	return
}

func (s *StringSet) Remove(value string) {
	if !s.Contains(value) {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.rm[value]

	delete(s.rm, value)
	delete(s.m, idx)
	s.k = append(s.k[:idx], s.k[idx+1:]...)
}

func (s *StringSet) Contains(value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.rm[value]
	return ok
}

func (s *StringSet) List() []string {
	l := len(s.k)
	values := make([]string, l)

	for i := range s.k {
		values[l-i-1] = s.m[i]
	}

	return values
}
