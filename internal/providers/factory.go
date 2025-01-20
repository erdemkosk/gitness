package providers

import "fmt"

type ProviderFactory struct {
	providers map[string]func(config map[string]string) (CommitProvider, error)
}

func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		providers: make(map[string]func(config map[string]string) (CommitProvider, error)),
	}

	factory.Register("github", func(config map[string]string) (CommitProvider, error) {
		token, exists := config["token"]
		if !exists {
			return nil, fmt.Errorf("GITHUB_TOKEN is required")
		}
		return NewGitHubProvider(token), nil
	})

	factory.Register("bitbucket", func(config map[string]string) (CommitProvider, error) {
		clientID, exists := config["clientID"]
		if !exists {
			return nil, fmt.Errorf("BITBUCKET_CLIENT_ID is required")
		}
		clientSecret, exists := config["clientSecret"]
		if !exists {
			return nil, fmt.Errorf("BITBUCKET_CLIENT_SECRET is required")
		}

		provider, err := NewBitbucketProvider(clientID, clientSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to create Bitbucket provider: %v", err)
		}
		return provider, nil
	})

	return factory
}

func (f *ProviderFactory) Register(name string, creator func(config map[string]string) (CommitProvider, error)) {
	f.providers[name] = creator
}

func (f *ProviderFactory) CreateProvider(providerType string, config map[string]string) (CommitProvider, error) {
	creator, exists := f.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
	return creator(config)
}
