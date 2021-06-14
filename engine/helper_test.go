package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeStepMap(t *testing.T) {
	steps := []*Step{
		{
			Name: "step1",
		},
		{
			Name:     "step2",
			Requires: []string{"step1"},
		},
	}

	got := MakeStepMap(steps)

	assert.Equal(t, got["step1"], steps[0])
	assert.Equal(t, got["step2"], steps[1])
}

func TestMakeGraph(t *testing.T) {
	steps := []*Step{
		{
			Name: "step1",
		},
		{
			Name:     "step2",
			Requires: []string{"step1", "step4"},
		},
		{
			Name:     "step3",
			Requires: []string{"step1", "step2"},
		},
		{
			Name:     "step4",
			Requires: []string{"step2"},
		},
	}

	got := MakeGraph(steps)

	assert.Equal(t, got["step1"], []string{})
	assert.Equal(t, got["step2"], []string{"step1", "step4"})
	assert.Equal(t, got["step3"], []string{"step1", "step2"})
	assert.Equal(t, got["step4"], []string{"step2"})
}

func TestMakeDfs(t *testing.T) {
	steps := []*Step{
		{
			Name: "step1",
		},
		{
			Name:     "step2",
			Requires: []string{"step1", "step4"},
		},
		{
			Name:     "step3",
			Requires: []string{"step1", "step2"},
		},
		{
			Name:     "step4",
			Requires: []string{"step2"},
		},
	}

	graph := MakeGraph(steps)

	got := Dfs(graph, "step2")

	assert.Equal(t, got, []string{"step1", "step4", "step2"})

	got = Dfs(graph, "step3")
	assert.Equal(t, got, []string{"step1", "step4", "step2", "step3"})
}
