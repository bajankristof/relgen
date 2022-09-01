package internal

import (
	"github.com/bajankristof/relgen/internal/conventionalcommits"
	"github.com/bajankristof/relgen/internal/semver"
	"time"
)

type Release struct {
	bump      string
	Version   *semver.Version `json:"version"`
	Changelog Changelog       `json:"changelog"`
	Date      time.Time       `json:"date"`
}

func NewRelease(version *semver.Version) *Release {
	return &Release{
		bump:      semver.NONE,
		Version:   semver.SelectLatest(semver.NewEmptyVersion(), version),
		Changelog: Changelog{},
		Date:      time.Now(),
	}
}

func (rel *Release) Push(cc *conventionalcommits.ConventionalCommit, spec *ChangeSpec) *Release {
	rel.Changelog[spec.Category] = append(rel.Changelog[spec.Category], cc)

	if cc.IsBreakingChange() {
		rel.bump = semver.MAJOR
	} else {
		rel.bump = semver.SelectGreaterBumpSpec(rel.bump, spec.Bump)
	}

	return rel
}

func (rel *Release) Close(tag string, metadata string) *Release {
	rel.Version.BumpWithSpec(rel.bump)
	rel.Version.WithPreReleaseTag(tag)
	rel.Version.Metadata = metadata
	rel.bump = semver.NONE
	return rel
}
