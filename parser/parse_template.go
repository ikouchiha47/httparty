package parser

import (
	"regexp"
	"strings"
)

func GetParsedString(str string) (map[string][]string, bool) {
	rx := regexp.MustCompile(`{{([^{{]+)}}`)
	res := rx.FindAllStringSubmatch(str, -1)

	if len(res) == 0 {
		return nil, false
	}

	result := map[string][]string{}

	for _, match := range res {
		matched := match[1]
		parts := strings.Split(matched, "::")

		switch len(parts) {
		case 1:
			result[match[0]] = []string{"", "", parts[0]}
		case 2:
			// e.g. "$config::client_id" -> prefix=$config, no resKey group, key=client_id
			result[match[0]] = []string{parts[0], "", parts[1]}
		default:
			result[match[0]] = parts
		}
	}

	return result, true
}
