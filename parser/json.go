package parser

import (
	"regexp"
	"strings"
)

var jsonTagRegexp = regexp.MustCompile(`json:"(.*?)"`)

// getJsonTag 获取标签中关于json部分的定义
func getJsonTag(tag string) string {
	matches := jsonTagRegexp.FindStringSubmatch(tag)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

// parseJsonTag splits a struct field's json tag into its name and
// comma-separated options.
func parseJsonTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, ""
}

// jsonTagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
