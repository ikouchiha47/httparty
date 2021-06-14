package engine

import (
	"github.com/stretchr/testify/require"
	"httparty/collections"
	"testing"
)

func TestStep_RenderHeadersAndBody(t *testing.T) {
	t.Run("parse non response step", func(t *testing.T) {
		step := &Step{
			Name: "step1",
			Request: Request{
				Headers: Headers{"client_id": "{{$config::client_id}}"},
			},
		}

		store := collections.NewOrderedMap()
		store.Add("client_id", "ClientID")

		st := step.RenderHeadersAndBody(store)
		require.Equal(t, st.Request.Headers, Headers{"client_id": "ClientID"})
	})

	t.Run("parse a response step", func(t *testing.T) {
		step := &Step{
			Name:     "step2",
			Request:  Request{
				Headers: Headers{"token": "Bearer {{$resp::access_token::token}}"},
			},
		}

		store := collections.NewOrderedMap()
		store.Add("$resp::access_token", Response{Body: Body{"token": "AccessToken"}})

		st := step.RenderHeadersAndBody(store)
		require.Equal(t, st.Request.Headers, Headers{"token": "Bearer AccessToken"})
	})
}
