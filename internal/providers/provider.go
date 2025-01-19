package providers

// CommitProvider defines the interface for different VCS providers
type CommitProvider interface {
	FetchCommits(owner, repo string) (map[string]int, error)
}
