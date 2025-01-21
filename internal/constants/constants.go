package constants

// Bus Factor thresholds
const (
	BusFactorCriticalThreshold      = 2
	BusFactorWarningThreshold       = 4
	BusFactorContributionPercentage = 80.0
)

// Knowledge Distribution thresholds
const (
	KnowledgeScoreCriticalThreshold = 25.0
	KnowledgeScoreWarningThreshold  = 50.0
	KnowledgeScoreMaximum           = 100.0
)

// Contributor Activity thresholds
const (
	ContributorActivityCriticalThreshold = 30.0
	ContributorActivityWarningThreshold  = 50.0
	SignificantContributionPercentage    = 1.0
)

// Recent Activity thresholds
const (
	RecentContributorsCriticalThreshold = 2
	RecentContributorsWarningThreshold  = 5
	RecentActivityMonths                = 3
)

// API Rate Limits
const (
	MaxConcurrentRequests = 10
	RequestRateLimit      = 200 // milliseconds
	MaxPageSize           = 100
	MaxPages              = 30
)

// Output Formatting
const (
	MaxDisplayedContributors = 5
	DateFormat               = "2006-01-02"
	TimeFormat               = "15:04"
)

// Duration Formats
const (
	DurationDay   = "d"
	DurationWeek  = "w"
	DurationMonth = "m"
	DurationYear  = "y"
)

// Default Values
const (
	DefaultDuration     = "6m"
	DefaultOutputFormat = "console"
	DefaultBatchSize    = 1000
	DefaultBufferSize   = 5000
)

// Error Messages
const (
	ErrEmptyOwnerRepo       = "owner and repo cannot be empty"
	ErrNoCommitsFound       = "no commits found in repository"
	ErrFailedToFetchCommits = "failed to fetch commits"
	ErrInvalidDuration      = "invalid duration format"
	ErrInvalidOutputFormat  = "invalid output format"
)

// HTTP Related
const (
	DefaultTimeout = 30 // seconds
	MaxRetries     = 3
	RetryDelay     = 1000 // milliseconds
)

// Cache Related
const (
	CacheExpiration = 24   // hours
	MaxCacheSize    = 1000 // items
)

// File Paths
const (
	DefaultConfigPath = ".gitness.yaml"
	DefaultOutputPath = "gitness_report"
)
