package mocking

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type MockRepository struct {
	LogReturn  *MockCommitIter
	LogCalls   []*git.LogOptions
	TagsReturn *MockReferenceIter
}

type MockReferenceIter struct {
	Error      error
	References []*plumbing.Reference
}

type MockCommitIter struct {
	Error   error
	Commits []*object.Commit
}

func (repo *MockRepository) Log(options *git.LogOptions) (object.CommitIter, error) {
	if repo.LogReturn.Error != nil {
		return nil, repo.LogReturn.Error
	}

	repo.LogCalls = append(repo.LogCalls, options)
	return repo.LogReturn, nil
}

func (repo *MockRepository) Tags() (storer.ReferenceIter, error) {
	if repo.TagsReturn.Error != nil {
		return nil, repo.TagsReturn.Error
	}

	return repo.TagsReturn, nil
}

func (iter *MockReferenceIter) Next() (*plumbing.Reference, error) {
	return nil, nil
}

func (iter *MockReferenceIter) ForEach(callback func(reference *plumbing.Reference) error) error {
	for _, ref := range iter.References {
		err := callback(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (iter *MockReferenceIter) Close() {
}

func (iter *MockCommitIter) Next() (*object.Commit, error) {
	return nil, nil
}

func (iter *MockCommitIter) ForEach(callback func(commit *object.Commit) error) error {
	for _, commit := range iter.Commits {
		err := callback(commit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (iter *MockCommitIter) Close() {
}
