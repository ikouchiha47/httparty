// A real, runnable demo of httparty: a local mock server plays the
// role of an OAuth-protected API, and rules.json describes a two-step
// fleet that (1) fetches a token using a pre-seeded config value, then
// (2) uses that token's response to authenticate a second request.
//
// Run it: go run ./example
package main

import (
	"encoding/json"
	"fmt"
	"httparty"
	"httparty/collections"
	"httparty/engine"
	"log"
	"net/http"
)

func startMockServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Client-Secret") != "s3cr3t-demo-key" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "bad client secret"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": "tok_abc123"})
	})

	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok_abc123" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "bad or missing token"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"name": "httparty demo user",
			"plan": "free",
		})
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:19191", mux))
}

func main() {
	go startMockServer()

	fleets, err := engine.ParseFile("./example/rules.json")
	if err != nil {
		log.Fatalf("failed to parse rules.json: %v", err)
	}

	// pre-seed the store: this is what {{$config::client_secret}} in
	// get_token's request resolves against, not a previous step's
	// response (that's what $resp:: is for).
	store := collections.NewOrderedMap()
	store.Add("client_secret", "s3cr3t-demo-key")

	hp := httparty.NewHttparty(httparty.WithStore(store))

	result, err := hp.RunIt(fleets[0])
	if err != nil {
		log.Fatalf("fleet failed: %v", err)
	}

	tokenResp := result.GetRespValue("get_token")
	profileResp := result.GetRespValue("get_profile")

	fmt.Println("get_token response: ", tokenResp)
	fmt.Println("get_profile response:", profileResp)
}
