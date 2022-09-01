package conventionalcommits

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"testing"
)

func TestNewConventionalCommit(t *testing.T) {
	commit := &object.Commit{Message: `testing(foo)!: lorem ipsum

dolor sit amet

My-Funky-Footer #ok
BREAKING CHANGE: it kinda broke
wHy-nOt: true

`}

	cc, err := NewConventionalCommit(commit)
	switch true {
	case err != nil:
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected error to be <nil>, got %v`, commit, cc, err, err)
	case cc.Type != "testing":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected type to be "%s", got "%s"`, commit, cc, err, "testing", cc.Type)
	case cc.Scope != "foo":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected scope to be "%s", got "%s"`, commit, cc, err, "foo", cc.Scope)
	case !cc.Exclamation:
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected exclamation to be true, got false`, commit, cc, err)
	case cc.Description != "lorem ipsum":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected description to be "%s", got "%s"`, commit, cc, err, "lorem ipsum", cc.Description)
	case cc.Body != "dolor sit amet":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected body to be "%s", got "%s"`, commit, cc, err, "dolor sit amet", cc.Body)
	case cc.Footers == nil:
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected footers to NOT be <nil>`, commit, cc, err)
	case cc.Footers["my-funky-footer"] != "ok":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected footers/my-funky-footer to be "%s", got "%s"`, commit, cc, err, "ok", cc.Footers["my-funky-footer"])
	case cc.Footers["why-not"] != "true":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected footers/why-not to be "%s", got "%s"`, commit, cc, err, "true", cc.Footers["why-not"])
	case cc.Footers["BREAKING CHANGE"] != "it kinda broke":
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected footers/BREAKING CHANGE to be "%s", got "%s"`, commit, cc, err, "it kinda broke", cc.Footers["BREAKING CHANGE"])
	}

	commit = &object.Commit{Message: `looks : almost good`}
	cc, err = NewConventionalCommit(commit)
	switch true {
	case err == nil:
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected error to NOT be <nil>`, commit, cc, err)
	case cc != nil:
		t.Fatalf(`NewConventionalCommit(%v) = (%v, %v), expected conventional commit to be <nil>`, commit, cc, err)
	}
}

func TestConventionalCommit_IsBreakingChange(t *testing.T) {
	var cc *ConventionalCommit
	cc, _ = NewConventionalCommit(&object.Commit{Message: "feat!: breaking"})
	if !cc.IsBreakingChange() {
		t.Fatalf(`(*ConventionalCommit(%v)).IsBreakingChange(), expected true, got false`, cc)
	}

	cc, _ = NewConventionalCommit(&object.Commit{Message: `feat: non-breaking

BREAKING CHANGE: SIKE`})
	if !cc.IsBreakingChange() {
		t.Fatalf(`(*ConventionalCommit(%v)).IsBreakingChange(), expected true, got false`, cc)
	}

	cc, _ = NewConventionalCommit(&object.Commit{Message: `feat: truly non-breaking

BREAKING-CHANGE: goteem again`})
	if !cc.IsBreakingChange() {
		t.Fatalf(`(*ConventionalCommit(%v)).IsBreakingChange(), expected true, got false`, cc)
	}

	cc, _ = NewConventionalCommit(&object.Commit{Message: "feat: truly truly non-breaking"})
	if cc.IsBreakingChange() {
		t.Fatalf(`(*ConventionalCommit(%v)).IsBreakingChange(), expected false, got true`, cc)
	}
}
