package injection

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Repository interface {
	Log(options *git.LogOptions) (object.CommitIter, error)
	Tags() (storer.ReferenceIter, error)
}

type NonMergeCommitIter struct {
	cache    map[plumbing.Hash]bool
	MaxDepth uint
}

func (iter *NonMergeCommitIter) ForEach(commits object.CommitIter, callback func(commit *object.Commit) error) error {
	iter.cache = map[plumbing.Hash]bool{}
	return iter.deepForEach(commits, callback, 0)
}

func (iter *NonMergeCommitIter) deepForEach(commits object.CommitIter, callback func(commit *object.Commit) error, depth uint) error {
	return commits.ForEach(func(commit *object.Commit) error {
		if iter.cache[commit.Hash] {
			return nil
		}

		iter.cache[commit.Hash] = true

		if commit.NumParents() < 2 {
			return callback(commit)
		}

		if depth < iter.MaxDepth {
			return iter.deepForEach(commit.Parents(), callback, depth+1)
		}

		return nil
	})
}
