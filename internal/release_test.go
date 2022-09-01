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
	case rel.Version != vsn:
		t.Fatalf(`NewRelease(%v) = %v, expected version to be %v, got %v`, vsn, rel, vsn, rel.Version)
	}
}

func TestRelease_Push(t *testing.T) {
	cc := &conventionalcommits.ConventionalCommit{}
	spec := &ChangeSpec{Bump: semver.MINOR, Category: "Tests"}
	rel := NewRelease(nil).Push(cc, spec)

	switch true {
	case rel.bump != semver.MINOR:
		t.Fatalf(`(*Release(%v)).Push(...), expected to set bump to "%s", got "%s"`, rel, semver.MINOR, rel.bump)
	case rel.Changelog[spec.Category][0] != cc:
		got := rel.Changelog[spec.Category][0]
		t.Fatalf(`(*Release(%v)).Push(...), expected to add %v to the changelog, got %v`, rel, cc, got)
	}

	cc = &conventionalcommits.ConventionalCommit{Exclamation: true}
	rel.Push(cc, spec)
	switch true {
	case rel.bump != semver.MAJOR:
		t.Fatalf(`(*Release(%v)).Push(...), expected to set bump to "%s", got "%s"`, rel, semver.MAJOR, rel.bump)
	case rel.Changelog[spec.Category][1] != cc:
		got := rel.Changelog[spec.Category][1]
		t.Fatalf(`(*Release(%v)).Push(...), expected to add %v to the changelog, got %v`, rel, cc, got)
	}
}

func TestRelease_Close(t *testing.T) {
	vsn, _ := semver.NewVersion("1.0.0")
	rel := NewRelease(vsn)
	rel.bump = semver.MAJOR

	rel.Close("", "")
	switch true {
	case rel.Version.String() != "2.0.0":
		t.Fatalf(`(*Release(%v)).Close("", ""), expected to bump version to 2.0.0, got %v`, rel, rel.Version)
	case rel.bump != semver.NONE:
		t.Fatalf(`(*Release(%v)).Close("", ""), expected to reset bump, got "%s"`, rel, rel.bump)
	}
}
