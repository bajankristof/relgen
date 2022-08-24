package internal

import (
	"github.com/bajankristof/relgen/internal/conventionalcommits"
	"github.com/bajankristof/relgen/internal/semver"
	"testing"
)

func TestNewRelease(t *testing.T) {
	vsn, _ := semver.NewVersion("1.0.0")
	rel := NewRelease(vsn)

	switch true {
	case rel.bump != semver.NONE:
		t.Fatalf(`NewRelease(%v) = %v, expected bump to be "%s", got "%s"`, vsn, rel, semver.NONE, rel.bump)
	case rel.version != vsn:
		t.Fatalf(`NewRelease(%v) = %v, expected version to be %v, got %v`, vsn, rel, vsn, rel.version)
	}
}

func TestRelease_Push(t *testing.T) {
	cc := &conventionalcommits.ConventionalCommit{}
	spec := &ChangeSpec{Bump: semver.MINOR, Category: "Tests"}
	rel := NewRelease(nil).Push(cc, spec)

	switch true {
	case rel.bump != semver.MINOR:
		t.Fatalf(`(*Release(%v)).Push(...), expected to set bump to "%s", got "%s"`, rel, semver.MINOR, rel.bump)
	case rel.changelog[spec.Category][0] != cc:
		got := rel.changelog[spec.Category][0]
		t.Fatalf(`(*Release(%v)).Push(...), expected to add %v to the changelog, got %v`, rel, cc, got)
	}

	rel.closed = true
	cc = &conventionalcommits.ConventionalCommit{Exclamation: true}
	rel.Push(cc, spec)
	switch true {
	case rel.bump != semver.MINOR:
		t.Fatalf(`(*Release(%v)).Push(...), expected to keep bump unchanged (closed)`, rel)
	case len(rel.changelog[spec.Category]) != 1:
		t.Fatalf(`(*Release(%v)).Push(...), expected to keep changelog unchanged (closed)`, rel)
	}

	rel.closed = false
	rel.Push(cc, spec)
	switch true {
	case rel.bump != semver.MAJOR:
		t.Fatalf(`(*Release(%v)).Push(...), expected to set bump to "%s", got "%s"`, rel, semver.MAJOR, rel.bump)
	case rel.changelog[spec.Category][1] != cc:
		got := rel.changelog[spec.Category][1]
		t.Fatalf(`(*Release(%v)).Push(...), expected to add %v to the changelog, got %v`, rel, cc, got)
	}
}

func TestRelease_Close(t *testing.T) {
	vsn, _ := semver.NewVersion("1.0.0")
	rel := NewRelease(vsn)
	rel.bump = semver.MAJOR

	rel.Close("", "")
	switch true {
	case !rel.closed:
		t.Fatalf(`(*Release(%v)).Close("", ""), expected to close the release`, rel)
	case rel.version.String() != "2.0.0":
		t.Fatalf(`(*Release(%v)).Close("", ""), expected to bump version to 2.0.0, got %v`, rel, rel.version)
	}

	rel.Close("foo", "bar")
	if rel.version.String() != "2.0.0" {
		t.Fatalf(`(*Release(%v)).Close(...), expected to keep version (closed)`, rel)
	}
}

func TestRelease_MarshalJSON(t *testing.T) {
	rel := NewRelease(nil)
	rel.bump = semver.MAJOR
	rel.Close("", "foo")

	expect := `{"version":"0.1.0+foo","changelog":{}}`
	if bytes, _ := rel.MarshalJSON(); string(bytes) != expect {
		t.Fatalf(`(*Release(%v)).MarshalJSON(...), expected '%s', got '%s'`, rel, expect, string(bytes))
	}
}
