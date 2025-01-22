package constants

// MetricThresholds contains all metric-related thresholds
type MetricThresholds struct {
	BusFactor        BusFactorThresholds
	KnowledgeScore   KnowledgeScoreThresholds
	ContributorStats ContributorStatsThresholds
}

// BusFactorThresholds defines thresholds for bus factor analysis
type BusFactorThresholds struct {
	Critical   int
	Warning    int
	Percentage float64
}

// KnowledgeScoreThresholds defines thresholds for knowledge distribution
type KnowledgeScoreThresholds struct {
	Critical float64
	Warning  float64
	Maximum  float64
}

// ContributorStatsThresholds defines thresholds for contributor statistics
type ContributorStatsThresholds struct {
	ActivityCritical   float64
	ActivityWarning    float64
	SignificantContrib float64
	RecentCritical     int
	RecentWarning      int
	RecentMonths       int
}

var Metrics = MetricThresholds{
	BusFactor: BusFactorThresholds{
		Critical:   2,
		Warning:    4,
		Percentage: 80.0,
	},
	KnowledgeScore: KnowledgeScoreThresholds{
		Critical: 25.0,
		Warning:  50.0,
		Maximum:  100.0,
	},
	ContributorStats: ContributorStatsThresholds{
		ActivityCritical:   30.0,
		ActivityWarning:    50.0,
		SignificantContrib: 1.0,
		RecentCritical:     2,
		RecentWarning:      5,
		RecentMonths:       3,
	},
}

// APIConfig contains API-related constants
type APIConfig struct {
	MaxConcurrentRequests int
	RequestRateLimit      int
	MaxPageSize           int
	MaxPages              int
}

var API = APIConfig{
	MaxConcurrentRequests: 10,
	RequestRateLimit:      200,
	MaxPageSize:           100,
	MaxPages:              30,
}

// OutputConfig contains output formatting related constants
type OutputConfig struct {
	MaxDisplayedContributors int
	DateFormat               string
	TimeFormat               string
}

var Output = OutputConfig{
	MaxDisplayedContributors: 5,
	DateFormat:               "2006-01-02",
	TimeFormat:               "15:04",
}

// DurationFormats contains duration format constants
var DurationFormats = struct {
	Day   string
	Week  string
	Month string
	Year  string
}{
	Day:   "d",
	Week:  "w",
	Month: "m",
	Year:  "y",
}

// Defaults contains default values
var Defaults = struct {
	Duration     string
	OutputFormat string
	BatchSize    int
	BufferSize   int
}{
	Duration:     "6m",
	OutputFormat: "console",
	BatchSize:    1000,
	BufferSize:   5000,
}

// ErrorMessages contains error message constants
var ErrorMessages = struct {
	EmptyOwnerRepo       string
	NoCommitsFound       string
	FailedToFetchCommits string
	InvalidDuration      string
	InvalidOutputFormat  string
}{
	EmptyOwnerRepo:       "owner and repo cannot be empty",
	NoCommitsFound:       "no commits found in repository",
	FailedToFetchCommits: "failed to fetch commits",
	InvalidDuration:      "invalid duration format",
	InvalidOutputFormat:  "invalid output format",
}

// HTTPConfig contains HTTP related constants
type HTTPConfig struct {
	DefaultTimeout int
	MaxRetries     int
	RetryDelay     int
}

var HTTP = HTTPConfig{
	DefaultTimeout: 30,
	MaxRetries:     3,
	RetryDelay:     1000,
}

// CacheConfig contains cache related constants
type CacheConfig struct {
	ExpirationHours int
	MaxSize         int
}

var Cache = CacheConfig{
	ExpirationHours: 24,
	MaxSize:         1000,
}

// FileConfig contains file path related constants
type FileConfig struct {
	DefaultConfigPath string
	DefaultOutputPath string
}

var Files = FileConfig{
	DefaultConfigPath: ".gitness.yaml",
	DefaultOutputPath: "gitness_report",
}
