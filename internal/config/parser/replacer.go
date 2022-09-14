package parser

import (
	"strings"
)

var (
	knownMap   map[string]string
	unknownMap map[string]bool
)

func init() {
	knownMap = make(map[string]string)
	unknownMap = make(map[string]bool)
	for _, v := range known {
		knownMap[v[0]] = v[1]
	}
}

type customReplacer struct{}

func (_ customReplacer) Replace(str string) string {
	// convert to known variable name
	if k, ok := knownMap[str]; ok {
		if k == "" {
			return strings.ToUpper(strings.ReplaceAll(str, ".", "_"))
		}
		return k
	}
	// otherwise replace automatically but warn the user
	unknownMap[str] = true
	return strings.ToUpper(strings.ReplaceAll(str, ".", "_"))
}

func UnknownEnvVariables() []string {
	result := make([]string, 0, len(unknownMap))
	for k := range unknownMap {
		result = append(result, k)
	}
	return result
}
