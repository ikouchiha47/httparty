package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetParsedString(t *testing.T) {
	t.Run("with proper template", func(t *testing.T) {
		resp, ok := GetParsedString("Bearer {{$resp::access_token::token}}")
		require.True(t, ok)

		wanted := map[string][]string{"{{$resp::access_token::token}}": []string{"$resp", "access_token", "token"}}
		assert.Equal(t, resp, wanted)
	})

	t.Run("with no templating", func(t *testing.T) {
		resp, ok := GetParsedString("Bearer $resp::access_token::token")
		require.False(t, ok)
		assert.Nil(t, resp)
	})

	t.Run("with template but no key parser", func(t *testing.T) {
		resp, ok := GetParsedString("Bearer {{resp_access_token}}")
		require.True(t, ok)

		want := map[string][]string{"{{resp_access_token}}": []string{"", "", "resp_access_token"}}
		assert.Equal(t, resp, want)
	})

	t.Run("with multiple templates", func(t *testing.T) {
		resp, ok := GetParsedString("Bearer {{$r::a::b.c}} {{$r::c::d}} Test{{$x::y::z}} What{{ever}}")
		require.True(t, ok)

		want := map[string][]string{
			"{{$r::a::b.c}}": {"$r", "a", "b.c"},
			"{{$r::c::d}}":   {"$r", "c", "d"},
			"{{$x::y::z}}":   {"$x", "y", "z"},
			"{{ever}}":       {"", "", "ever"},
		}

		assert.Equal(t, resp, want)
	})
}

