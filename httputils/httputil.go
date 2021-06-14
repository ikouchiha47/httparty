package httputils

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func MakeFormEncodedBody(requestBody map[string]interface{}) string {
	data := url.Values{}
	for key, value := range requestBody {
		data.Set(key, fmt.Sprintf("%v", value))
	}

	return data.Encode()
}

func MakeJSONBody(requestBody map[string]interface{}) (string, error) {
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
