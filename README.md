## httparty

This is a rule engine that parses a json file which contains a series of json requests to be made
and returns the final response.


Each `Step` is an http request, represented as such:

```
{
    "name": "access_token",
    "request": {
        "url": "http://localhost:8080/v1/oauth/accesstoken",
        "method": "POST",
        "type": "application/x-www-form-urlencoded",
        "accept": "application/json",
        "headers": {
            "Authorization": "Basic c2FtcGxldXNlcjpzYW1wbGVwYXNzd29yZA==",
            "apiKey": "2Dk"
        },
        "body": {
            "grant_type": "client_credentials",
            "scope": "public"
        },
        "requires": []
    }
}
```

- `name` is an unique identifier
- `request` object contains the url, method, content-type, accept type, headers and body
- `requires` contains the unique names of the rules that are required to make this request

The `header` and `body` can also take values depending on the preceding request mentioned in `requires` array

Example:

```
{
    "exitpoint": 0,
    "steps": [
        {
            "name": "rakshas_token",
            "request": {
                "url": "",
                "method": "",
                "type": "",
                "accept": "",
                "headers": {
                    "Authorization": "Bearer {{$resp::access_token::token}}"
                }
            },
            "requires": ["access_token"]
        },
       {
            "name": "access_token",
            "request": {
                "url": "",
                "method": "",
                "type": "",
                "accept": ""
            }
        }
    ]
}
```
**All the reponse are stored under "$resp" key**


Here:

- `exitpoint` refers to the final request that has to be made (should use a better name)
- `{{$resp::access_token::token}}` all values separated by `::` signifies:
    - `$resp` -> that a value from response needs to be used (could be `$config` or anything to group the store)
    - `access_token` -> that the response from the request in rule named `access_token` needs to be used
    - `token` -> is the key in the response body (this could use something like a `jq` style nested json like `.results[0].name`)
    


## Live example

`example/` is a real, runnable demo, not a mock at the unit-test level.
It spins up an actual local HTTP server standing in for an
OAuth-protected API, and runs a real two-step fleet against it over the
network:

1. `get_token` POSTs to `/oauth/token`, authenticating with a client
   secret pulled from a pre-seeded store via `{{$config::client_secret}}`
   (not a previous step's response - that's what `$resp::` is for).
2. `get_profile` GETs `/profile`, using the token from step 1's response
   via `{{$resp::get_token::token}}`. `requires: ["get_token"]` is what
   makes the engine run step 1 first.

Run it:

```sh
go run ./example
```

Expected output:

```
get_token response:  {map[token:tok_abc123]}
get_profile response: {map[name:httparty demo user plan:free]}
```

See `example/main.go` for the mock server and `example/rules.json` for
the fleet definition.

Thats all folks!
