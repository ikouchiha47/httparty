package collections

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestOrderedStringSet(t *testing.T) {
	ordset := NewStringSet()

	values := []string{"abcd", "ijkl", "efgh", "ijkl"}

	for _, value := range values {
		ordset.Add(value)
	}

	expected := []string{"abcd", "ijkl", "efgh"}

	got := ordset.List()

	assert.Equal(t, expected, got)
}
