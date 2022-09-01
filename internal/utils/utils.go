package utils

import (
	"regexp"
)

type NamedRegexpGroupIter struct {
	*regexp.Regexp
}

func (iter *NamedRegexpGroupIter) ForEach(subject string, callback func(group string, match string)) bool {
	match := iter.Regexp.FindStringSubmatch(subject)
	if len(match) < 1 {
		return false
	}

	for i, group := range iter.Regexp.SubexpNames() {
		if i == 0 || group == "" {
			continue
		}

		callback(group, match[i])
	}

	return true
}
