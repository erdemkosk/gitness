package providers

import (
	"time"
)

// CommitProvider defines the interface for different VCS providers
type CommitProvider interface {
	FetchCommits(owner, repo string) (map[string]CommitInfo, error)
}

type CommitInfo struct {
	Count      int
	LastCommit time.Time
}
