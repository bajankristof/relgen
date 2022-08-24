package injection

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Repository interface {
	Log(options *git.LogOptions) (object.CommitIter, error)
	Tags() (storer.ReferenceIter, error)
}
