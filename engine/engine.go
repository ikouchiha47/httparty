package engine

import (
	"encoding/json"
	"fmt"
	"httparty/collections"
	"httparty/errorrs"
	"httparty/parser"
	"io/ioutil"
	"strings"
)

type Step struct {
	Name     string   `json:"name"`
	Request  Request  `json:"request"`
	Requires []string `json:"requires"`
}

type Fleet struct {
	Steps     []*Step `json:"steps,omitempty"`
	ExitPoint int     `json:"exitpoint,omitempty"`
}

type FleetI interface{}

func ParseFile(filePath string) ([]Fleet, error) {
	fleet := []Fleet{}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fleet, errorrs.InvalidFile
	}

	if err = json.Unmarshal(b, &fleet); err != nil {
		return fleet, errorrs.InvalidFormat
	}

	return fleet, nil
}

func (s *Step) RenderHeadersAndBody(store *collections.OrderedMap) *Step {
	rq, headers, body := s.Request, s.Request.Headers, s.Request.Body

	for key, value := range headers {
		headers[key] = s.getValue(value, store)
	}

	for key, value := range body {
		body[key] = s.getValue(fmt.Sprintf("%v", value), store)
	}

	return &Step{
		Name: s.Name,
		Request: Request{
			URL:     rq.URL,
			Method:  rq.Method,
			Type:    rq.Type,
			Accept:  rq.Accept,
			Headers: headers,
			Body:    body,
		},
	}
}

func (s *Step) getValue(templateStr string, store *collections.OrderedMap) string {
	parseMap, ok := parser.GetParsedString(templateStr)
	if !ok {
		return templateStr
	}

	for matched, values := range parseMap {
		prefix, resKey, key := values[0], values[1], values[2]
		if prefix == "" || resKey == "" {
			templateStr = strings.ReplaceAll(templateStr, matched, key)
			continue
		}

		var value string
		var ok bool

		if prefix == RespPrefix {
			value, ok = getValueForResp(store, fmt.Sprintf("%s::%s", prefix,resKey), key)
		} else {
			value, ok = getValue(store, key)
		}

		if !ok {
			continue
		}

		templateStr = strings.ReplaceAll(templateStr, matched, value)
	}

	return templateStr
}

func getValueForResp(store *collections.OrderedMap, resKey, key string) (string, bool) {
	kv, ok := store.Get(resKey)
	if !ok {
		return "", false
	}

	// TODO: write a better json dot parser
	value, ok := kv.(Response).Body[key]
	if !ok {
		return "", false
	}

	return value.(string), true
}

func getValue(store *collections.OrderedMap, key string) (string, bool) {
	value, ok := store.Get(key)
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%v", value), true
}
