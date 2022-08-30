package injection

import (
	"errors"
	"github.com/bajankristof/relgen/internal/mocking"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"testing"
)

var errForEach = errors.New("test")

func TestNonMergeCommitIter_ForEach(t *testing.T) {
	log := &mocking.MockCommitIter{Commits: []*object.Commit{
		{Message: "foo", Hash: plumbing.NewHash("123"), ParentHashes: []plumbing.Hash{}},
		{Message: "bar", Hash: plumbing.NewHash("456"), ParentHashes: []plumbing.Hash{}},
		{Message: "baz", Hash: plumbing.NewHash("789"), ParentHashes: []plumbing.Hash{}},
	}}

	log.Commits = append([]*object.Commit{log.Commits[1]}, log.Commits...)
	iter := &NonMergeCommitIter{MaxDepth: 10}

	var commits []*object.Commit
	err := iter.ForEach(log, func(commit *object.Commit) error {
		if commit.Hash == log.Commits[len(log.Commits)-1].Hash {
			return errForEach
		}

		commits = append(commits, commit)
		return nil
	})

	switch true {
	case err != errForEach:
		t.Fatalf("(*NonMergeCommitIter(%v)).ForEach(...) = %v, expected error to be %v, got %v", iter, err, errForEach, err)
	case len(commits) != 2:
		expect := 2
		got := len(commits)
		t.Fatalf("(*NonMergeCommitIter(%v)).ForEach(...) = %v, expected number of iterations to be %v, got %v", iter, err, expect, got)
	}
}

func TestNonMergeCommitIter_deepForEach(t *testing.T) {
	log := &mocking.MockCommitIter{Commits: []*object.Commit{
		{Message: "foo", Hash: plumbing.NewHash("123"), ParentHashes: []plumbing.Hash{plumbing.NewHash("111"), plumbing.NewHash("222")}},
		{Message: "bar", Hash: plumbing.NewHash("456"), ParentHashes: []plumbing.Hash{plumbing.NewHash("333"), plumbing.NewHash("444")}},
	}}

	iter := &NonMergeCommitIter{MaxDepth: 10}
	iter.cache = map[plumbing.Hash]bool{}

	var commits []*object.Commit
	iter.deepForEach(log, func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	}, 10)

	got := len(commits)
	if got > 0 {
		t.Fatalf("(*NonMergeCommitIter(%v)).deepForEach(...), expected number of iterations to be %v, got %v", iter, 0, got)
	}
}
