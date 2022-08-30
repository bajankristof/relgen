package internal

import (
	"errors"
	"github.com/bajankristof/relgen/internal/mocking"
	"github.com/bajankristof/relgen/internal/semver"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"testing"
)

func TestNewReleaseBuilder(t *testing.T) {
	repo := &mocking.MockRepository{}
	cfg := &Config{}
	builder := NewReleaseBuilder(repo, cfg)

	if builder.Repository != repo {
		t.Fatalf("NewReleaseBuilder(%v, %v) %v, expected repository to be %v, got %v", repo, cfg, builder, repo, builder.Repository)
	}

	if builder.Config != cfg {
		t.Fatalf("NewReleaseBuilder(%v, %v) %v, expected repository to be %v, got %v", repo, cfg, builder, cfg, builder.Config)
	}
}

func TestReleaseBuilder_Build(t *testing.T) {
	repo := &mocking.MockRepository{
		LogReturn: &mocking.MockCommitIter{},
		TagsReturn: &mocking.MockReferenceIter{
			References: []*plumbing.Reference{
				plumbing.NewReferenceFromStrings("2.0.0-test", "a"),
				plumbing.NewReferenceFromStrings("foo", "b"),
				plumbing.NewReferenceFromStrings("3.0.0-zod", "d"),
				plumbing.NewReferenceFromStrings("2.0.0-test.1", "c"),
				plumbing.NewReferenceFromStrings("1.1.0", "d"),
			},
		},
	}

	builder := NewReleaseBuilder(repo, &Config{PreRelease: "test"})
	rel, err := builder.Build()

	switch true {
	case err != nil:
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected error to be <nil>, got %v", builder, rel, err, err)
	case repo.LogCalls[0].From != repo.TagsReturn.References[1].Hash():
		got := repo.LogCalls[0].From
		expect := repo.TagsReturn.References[1].Hash()
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected to call log with %v, got %v", builder, rel, err, expect, got)
	}
}

func TestReleaseBuilder_BuildNoTags(t *testing.T) {
	repo := &mocking.MockRepository{
		TagsReturn: &mocking.MockReferenceIter{},
		LogReturn:  &mocking.MockCommitIter{},
	}

	builder := NewReleaseBuilder(repo, &Config{})
	rel, err := builder.Build()

	switch true {
	case err != nil:
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected error to be <nil>, got %v", builder, rel, err, err)
	case repo.LogCalls[0].From != plumbing.ZeroHash:
		got := repo.LogCalls[0].From
		expect := plumbing.ZeroHash
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected to call log with %v, got %v", builder, rel, err, expect, got)
	}
}

func TestReleaseBuilder_BuildTagsError(t *testing.T) {
	repo := &mocking.MockRepository{
		TagsReturn: &mocking.MockReferenceIter{Error: errors.New("nok")},
		LogReturn:  &mocking.MockCommitIter{},
	}

	builder := NewReleaseBuilder(repo, &Config{})
	rel, err := builder.Build()

	switch true {
	case err != repo.TagsReturn.Error:
		expect := repo.TagsReturn.Error
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected error to be %v, got %v", builder, rel, err, expect, err)
	case rel != nil:
		t.Fatalf("(*ReleaseBuilder(%v)).Build() = (%v, %v), expected release to be <nil>, got %v", builder, rel, err, rel)
	}
}

func TestReleaseBuilder_BuildSince(t *testing.T) {
	repo := &mocking.MockRepository{
		LogReturn: &mocking.MockCommitIter{Commits: []*object.Commit{
			{Message: "looks : almost good", Hash: plumbing.NewHash("123"), ParentHashes: []plumbing.Hash{}},
			{Message: "test: nice and conventional", Hash: plumbing.NewHash("456"), ParentHashes: []plumbing.Hash{}},
			{Message: "foo: conventional without purpose", Hash: plumbing.NewHash("789"), ParentHashes: []plumbing.Hash{}},
		}},
	}

	category := "Tests"
	builder := NewReleaseBuilder(repo, &Config{
		PreRelease:    "test",
		BuildMetadata: "foo",
		ChangeSpec: []ChangeSpec{
			{Type: &TypeSpec{regexp.MustCompile("^test$")}, Bump: semver.MAJOR, Category: category},
		},
	})

	rel, err := builder.BuildSince(nil)

	switch true {
	case err != nil:
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected error to be <nil>, got %v", builder, rel, err, err)
	case rel.version.String() != "0.1.0-test+foo":
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected release version to be 0.1.0-test+foo, got %v", builder, rel, err, rel.version)
	case len(rel.changelog[category]) != 1:
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected release changelog length to be 1, got %d", builder, rel, err, len(rel.changelog[category]))
	case rel.changelog[category][0].Commit != repo.LogReturn.Commits[1]:
		expect := repo.LogReturn.Commits[1]
		got := rel.changelog[category][0].Commit
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected release changelog to contain %v, got %v", builder, rel, err, expect, got)
	}
}

func TestReleaseBuilder_BuildSinceLogError(t *testing.T) {
	repo := &mocking.MockRepository{
		LogReturn: &mocking.MockCommitIter{
			Error: errors.New("nok"),
		},
	}

	builder := NewReleaseBuilder(repo, &Config{})
	rel, err := builder.BuildSince(nil)

	switch true {
	case err != repo.LogReturn.Error:
		expect := repo.LogReturn.Error
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected error to be %v, got %v", builder, rel, err, expect, err)
	case rel != nil:
		t.Fatalf("(*ReleaseBuilder(%v)).BuildSince(nil) = (%v, %v), expected release to be <nil>, got %v", builder, rel, err, rel)
	}
}
