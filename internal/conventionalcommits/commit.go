package conventionalcommits

import (
	"errors"
	"github.com/bajankristof/relgen/internal/utils"
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
	cc := &ConventionalCommit{Commit: commit, Footers: map[string]string{}}
	chunks := strings.Split(strings.TrimRight(commit.Message, "\n"), "\n")
	if !cc.parseMessage(chunks[0]) {
		return nil, errors.New("commit is not conventional")
	}

	cc.parseBodyAndFooters(chunks[1:])
	return cc, nil
}

func (cc *ConventionalCommit) AddFooter(key string, value string) {
	if key != "BREAKING CHANGE" && key != "BREAKING-CHANGE" {
		key = strings.ToLower(key)
	}

	cc.Footers[key] = value
}

func (cc *ConventionalCommit) HasFooter(key string) bool {
	_, ok := cc.Footers[key]
	return ok
}

func (cc *ConventionalCommit) IsBreakingChange() bool {
	if cc.Exclamation {
		return true
	}

	return cc.HasFooter("BREAKING CHANGE") ||
		cc.HasFooter("BREAKING-CHANGE")
}

func (cc *ConventionalCommit) parseMessage(message string) bool {
	iter := &utils.NamedRegexpGroupIter{Regexp: MessageRegex}
	return iter.ForEach(message, func(group string, match string) {
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
}

func (cc *ConventionalCommit) parseBodyAndFooters(chunks []string) {
	iter := &utils.NamedRegexpGroupIter{Regexp: FooterRegex}
	for i := len(chunks) - 1; i >= 0; i-- {
		var key, value string
		if !iter.ForEach(chunks[i], func(group string, match string) {
			switch group {
			case "key":
				key = match
			case "value":
				value = match
			}
		}) {
			break
		}

		cc.AddFooter(key, value)
		chunks = chunks[:i]
	}

	cc.Body = strings.TrimSpace(strings.Join(chunks, "\n"))
}
