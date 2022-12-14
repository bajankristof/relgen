package internal

import (
	"errors"
	"github.com/bajankristof/relgen/internal/conventionalcommits"
	"github.com/bajankristof/relgen/internal/injection"
	"github.com/bajankristof/relgen/internal/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var errBreak = errors.New("break")

type ReleaseBuilder struct {
	bump       string
	Repository injection.Repository
	Config     *Config
}

func NewReleaseBuilder(repository injection.Repository, config *Config) *ReleaseBuilder {
	return &ReleaseBuilder{Repository: repository, Config: config}
}

func (builder *ReleaseBuilder) Build() (*Release, error) {
	currentVsn, err := builder.ReadCurrentVersion()
	if err != nil {
		return nil, err
	}

	currentVsn = semver.SelectLatest(semver.NewEmptyVersion(), currentVsn)
	return builder.BuildSince(currentVsn)
}

func (builder *ReleaseBuilder) BuildSince(version *semver.Version) (*Release, error) {
	commits, err := builder.Repository.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}

	rel := NewRelease(builder.NewReleaseVersion(version))
	iter := injection.NonMergeCommitIter{MaxDepth: 1}
	err = iter.ForEach(commits, func(commit *object.Commit) error {
		if version != nil && version.IsReference(commit.Hash) {
			return errBreak
		}

		cc, err := conventionalcommits.NewConventionalCommit(commit)
		if err != nil {
			return nil
		}

		if cc.HasFooter("relgen-off") {
			return nil
		}

		_, spec := builder.Config.FindChangeSpec(cc)
		if spec == nil {
			return nil
		}

		rel.Push(cc, spec)
		return nil
	})

	if err != nil && err != errBreak {
		return nil, err
	}

	rel.Close(builder.Config.PreRelease, builder.Config.BuildMetadata)

	return rel, nil
}

func (builder *ReleaseBuilder) ReadCurrentVersion() (*semver.Version, error) {
	tags, err := builder.Repository.Tags()
	if err != nil {
		return nil, err
	}

	var vsn *semver.Version
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		tagVsn, err := semver.NewVersion(ref.Name().Short())
		if err != nil {
			return nil
		}

		if !tagVsn.MatchPreReleaseTag(builder.Config.PreRelease) {
			return nil
		}

		vsn = semver.SelectLatest(vsn, tagVsn.WithReference(ref))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return vsn, nil
}

func (builder *ReleaseBuilder) NewReleaseVersion(version *semver.Version) *semver.Version {
	vsn := &(*semver.SelectLatest(semver.NewEmptyVersion(), version))
	return vsn.WithPrefix(builder.Config.VersionPrefix)
}
