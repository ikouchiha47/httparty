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

		if len(parts) < 3 && len(parts) > 0 {
			result[match[0]] = []string{parts[0], parts[1], parts[1]}
		} else {
			result[match[0]] = parts
		}
	}

	return result, true
}
