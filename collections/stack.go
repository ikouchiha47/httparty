package collections

import "sync"

type StringStack struct {
	mu sync.Mutex
	s  []string
}

func NewStringStack(args ...string) *StringStack {
	s := []string{}

	for _, arg := range args {
		s = append(s, arg)
	}

	return &StringStack{s: s}
}

func (st *StringStack) Push(value string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.s = append(st.s, value)
}

func (st *StringStack) Pop() string {
	var v string

	st.mu.Lock()
	defer st.mu.Unlock()

	i := len(st.s) - 1

	v, st.s = st.s[i], st.s[:i]
	return v
}

func (st *StringStack) Empty() bool {
	return len(st.s) <= 0
}
