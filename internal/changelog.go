package internal

import (
	"github.com/bajankristof/relgen/internal/conventionalcommits"
)

type Changelog map[string][]*conventionalcommits.ConventionalCommit
