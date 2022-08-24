package conventionalcommits

import (
	"errors"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"
)

type ConventionalCommit struct {
	*object.Commit
	Type        string            `json:"type"`
	Scope       string            `json:"scope"`
	Exclamation bool              `json:"exclamation"`
	Description string            `json:"description"`
	Footers     map[string]string `json:"footers"`
	Body        string            `json:"body"`
}

var (
	MessageRegex = regexp.MustCompile("(?i)^(?P<type>[a-z]{2,})(\\((?P<scope>[a-z]+)\\))?(?P<exclamation>!)?: (?P<description>[^ ].*)$")
	FooterRegex  = regexp.MustCompile("(?i)^(?P<key>([a-z]+(-+[a-z]+)*|BREAKING[ -]CHANGE))((: | #)(?P<value>[^ ].*))?$")
)

func NewConventionalCommit(commit *object.Commit) (*ConventionalCommit, error) {
	cc := &ConventionalCommit{Commit: commit}

	lines := strings.Split(commit.Message, "\n")
	ok := matchNamedGroups(lines[0], MessageRegex, func(group string, match string) {
		switch group {
		case "type":
			cc.Type = match
		case "scope":
			cc.Scope = match
		case "exclamation":
			cc.Exclamation = match == "!"
		case "description":
			cc.Description = match
		}
	})

	if !ok {
		return nil, errors.New("commit is not conventional")
	}

	lines = lines[1:]
	for i := len(lines) - 1; i >= 0; i-- {
		var key, value string
		ok := matchNamedGroups(lines[i], FooterRegex, func(group string, match string) {
			switch group {
			case "key":
				key = match
			case "value":
				value = match
			}
		})

		if !ok {
			break
		}

		if key != "BREAKING CHANGE" && key != "BREAKING-CHANGE" {
			key = strings.ToLower(key)
		}

		if cc.Footers == nil {
			cc.Footers = map[string]string{}
		}

		cc.Footers[key] = value
		lines = lines[:i]
	}

	cc.Body = strings.TrimSpace(strings.Join(lines, "\n"))
	return cc, nil
}

func (cc *ConventionalCommit) IsBreakingChange() bool {
	if cc.Exclamation {
		return true
	}

	if cc.Footers == nil {
		return false
	}

	_, ok := cc.Footers["BREAKING CHANGE"]
	if ok {
		return true
	}

	_, ok = cc.Footers["BREAKING-CHANGE"]
	return ok
}

func matchNamedGroups(subject string, regex *regexp.Regexp, callback func(group string, match string)) bool {
	match := regex.FindStringSubmatch(subject)
	if len(match) < 1 {
		return false
	}

	for i, group := range regex.SubexpNames() {
		if i == 0 || group == "" {
			continue
		}

		callback(group, match[i])
	}

	return true
}
