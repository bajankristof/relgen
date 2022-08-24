package internal

import (
	"encoding/json"
	"github.com/bajankristof/relgen/internal/conventionalcommits"
	"github.com/bajankristof/relgen/internal/semver"
)

type Release struct {
	bump      string
	closed    bool
	version   *semver.Version
	changelog Changelog
}

func NewRelease(version *semver.Version) *Release {
	return &Release{
		bump:      semver.NONE,
		version:   semver.SelectLatest(semver.NewEmptyVersion(), version),
		changelog: Changelog{},
	}
}

func (rel *Release) Push(cc *conventionalcommits.ConventionalCommit, spec *ChangeSpec) *Release {
	if rel.closed {
		return rel
	}

	rel.changelog[spec.Category] = append(rel.changelog[spec.Category], cc)

	if cc.IsBreakingChange() {
		rel.bump = semver.MAJOR
	} else {
		rel.bump = semver.SelectGreaterBumpSpec(rel.bump, spec.Bump)
	}

	return rel
}

func (rel *Release) Close(tag string, metadata string) *Release {
	if rel.closed {
		return rel
	}

	rel.version.BumpWithSpec(rel.bump)
	rel.version.WithPreReleaseTag(tag)
	rel.version.Metadata = metadata
	rel.closed = true
	return rel
}

func (rel *Release) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version   *semver.Version `json:"version"`
		Changelog Changelog       `json:"changelog"`
	}{
		Version:   rel.version,
		Changelog: rel.changelog,
	})
}
