# Gitness

Gitness - Your repo's fitness witness! Track your bus factor before your code misses the bus.

![Gitness](/logo.png)

## Features

- Calculate repository bus factor and knowledge distribution
- Analyze contributor statistics and activity patterns
- Track recent contributor engagement
- Support for multiple VCS providers (GitHub, Bitbucket)
- Multiple output formats (Console, JSON, Markdown)
- Configurable analysis period (e.g., 6m, 1y, 30d)
- Branch-specific analysis support
- CI/CD pipeline integration
- Docker support

## Usage

```bash
# Analyze all time
gitness https://github.com/user/repo

# Analyze last 6 months
gitness --duration 6m https://github.com/user/repo

# Analyze last 1 year with JSON output
gitness --duration 1y --output json https://github.com/user/repo

# Analyze specific branch with Markdown output
gitness --branch feature-branch --output markdown https://github.com/user/repo

# Analyze specific branch for last 30 days
gitness --duration 30d --branch experimental --output console https://github.com/user/repo

# Analyze last week with console output (default)
gitness --duration 7d --output console https://github.com/user/repo

# Analyze last quarter with JSON output and specific branch
gitness --duration 3m --branch develop --output json https://github.com/user/repo

# Analyze last month with Markdown output and save to file
gitness --duration 1m --output markdown https://github.com/user/repo > report.md
```

## Metrics Explained

### Bus Factor 🚌
The "Bus Factor" represents the minimum number of developers that would need to be hit by a bus before the project is in serious trouble. It's calculated based on contributors who collectively account for 80% of contributions.

- 🔴 Critical (< 2): Project knowledge is dangerously concentrated
- 🟡 Warning (2-3): Limited knowledge distribution
- 🟢 Good (≥ 4): Healthy knowledge distribution

### Knowledge Distribution Score 📊
Measures how evenly the knowledge is distributed across all contributors (0-100).

- 🔴 Critical (< 25): Knowledge is heavily concentrated
- 🟡 Warning (25-50): Moderate knowledge concentration
- 🟢 Good (> 50): Well-distributed knowledge

### Active Contributor Ratio 👥
Percentage of contributors who have made significant contributions (>1% of total commits).

- 🔴 Critical (< 30%): Most contributors are occasional
- 🟡 Warning (30-50%): Moderate active participation
- 🟢 Good (> 50%): Healthy active participation

### Recent Contributors 🕒
Number of contributors active in last 3 months.

- 🔴 Critical (< 2): Low recent activity
- 🟡 Warning (2-4): Moderate recent activity
- 🟢 Good (≥ 5): High recent activity

### Commit Frequency Analysis ⏰
Detailed analysis of commit patterns and activity trends.

- 📅 Daily Average: Average number of commits per day
- 📅 Weekly Average: Average number of commits per week
- 📅 Monthly Average: Average number of commits per month
- 📅 Most Active Day: Day of the week with highest commit activity
- 🕒 Peak Activity Time: Hour of the day with most commits

## Example Output

```
    ██████╗ ██╗████████╗███╗   ██╗███████╗███████╗███████╗
   ██╔════╝ ██║╚══██╔══╝████╗  ██║██╔════╝██╔════╝██╔════╝
   ██║  ███╗██║   ██║   ██╔██╗ ██║█████╗  ███████╗███████╗
   ██║   ██║██║   ██║   ██║╚██╗██║██╔══╝  ╚════██║╚════██║
   ╚██████╔╝██║   ██║   ██║ ╚████║███████╗███████║███████║
    ╚═════╝ ╚═╝   ╚═╝   ╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝

Your repo's fitness witness! Track your bus factor before your code misses the bus.
═══════════════════════════════════════════════════════════════════════════════════

📊 Repository: user/repo
──────────────────────────────────────────────────

🌿 Branch: default
🕒 Analysis Period: Last 6 months

🎯 Core Metrics
────────────────
🚌 Bus Factor: 4
📚 Knowledge Distribution: 75.5%
📝 Total Commits: 330
👥 Active Contributors: 45.5%
🔄 Recent Contributors: 6

⏰ Commit Frequency
────────────────
📅 Daily Average: 3.77 commits
📅 Weekly Average: 26.38 commits
📅 Monthly Average: 113.07 commits
📅 Most Active Day: Monday
🕒 Peak Activity Time: 18:00

👥 Top Contributors
────────────────
👤 John Doe: 100 commits (30.3%)
👤 Jane Smith: 90 commits (27.3%)
👤 Bob Wilson: 80 commits (24.2%)
👤 Alice Brown: 60 commits (18.2%)
👤 Charlie Brown: 30 commits (9.1%)

... and 5 more contributors

═══════════════════════════════════════════════════════════════
                     Generated by Gitness                      
```

## Environment Variables

```bash
GITHUB_TOKEN=your_token
BITBUCKET_CLIENT_ID=your_client_id
BITBUCKET_CLIENT_SECRET=your_client_secret
COMMIT_HISTORY_DURATION=6m  # Optional: 6m, 1y, 30d etc.
REPOSITORY_BRANCH=main      # Optional: Specify branch to analyze
```

## Installation

### Using Go 
```bash
go install github.com/erdemkosk/gitness@latest
```
### Using Docker
```bash
docker build -t gitness .
docker run \                                                                                                                           ok 
  -e GITHUB_TOKEN="TOKEN" \
  -e REPOSITORY_URL="https://github.com/user/repo" \
  gitness --output json;
```

### Environment Variables

- `REPOSITORY_URL`: Target repository URL
- `GITHUB_TOKEN`: GitHub personal access token
- `BITBUCKET_CLIENT_ID`: Bitbucket OAuth client ID
- `BITBUCKET_CLIENT_SECRET`: Bitbucket OAuth client secret
- `OUTPUT_FORMAT`: Output format (console, json, markdown)
- `COMMIT_HISTORY_DURATION`: Analysis period (e.g., 6m, 1y, 30d)
- `REPOSITORY_BRANCH`: Branch to analyze (optional, defaults to repository's default branch)

### Repository URL Formats

#### GitHub
```
https://github.com/username/repository
```

#### Bitbucket
```
https://bitbucket.org/workspace/repository
```

### Authentication

#### GitHub
Generate a personal access token from GitHub Settings > Developer Settings > Personal Access Tokens.

#### Bitbucket
1. Go to Bitbucket Settings > OAuth Consumers
2. Create a new OAuth consumer
3. Use the Client ID and Client Secret in your environment variables

### CI/CD Integration

#### GitHub Actions

```yaml
name: Gitness Bus Factor Analysis

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  check-bus-factor:
    runs-on: ubuntu-latest
    steps:
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      
      - name: Run Bus Factor Analysis
        uses: docker://erdemkosk/gitness:latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITNESS_GITHUB_TOKEN }}
          REPOSITORY_URL: "https://github.com/${{ github.repository }}"
          OUTPUT_FORMAT: json
        id: analysis
```
#### GitHub Actions With Result

```yaml
name: Gitness Analysis

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Gitness Analysis
        id: gitness
        run: |
          OUTPUT=$(docker run \
            -e GITHUB_TOKEN="${{ secrets.GITNESS_GITHUB_TOKEN }}" \
            -e REPOSITORY_URL="https://github.com/${{ github.repository }}" \
            -e OUTPUT_FORMAT=markdown \
            erdemkosk/gitness:latest)
          echo "# 📊 Gitness Analysis Report" >> $GITHUB_STEP_SUMMARY
          echo "$OUTPUT" >> $GITHUB_STEP_SUMMARY
          echo "> This report is automatically generated by [Gitness](https://github.com/erdemkosk/gitness)" >> $GITHUB_STEP_SUMMARY
```

###### Example GitHub Action Result

[Example GitHub Action Run](https://github.com/erdemkosk/stock-exchange/actions/runs/12857735422)

## Architecture

- Clean architecture principles
- Strategy pattern for output formatting
- Factory pattern for VCS providers
- Dependency injection
- Environment-based configuration

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

### Command Line Arguments

- `--output`: Output format (console, json, markdown)
- `--duration`: Analysis period (e.g., 6m, 1y, 30d)
- `--branch`: Branch to analyze (optional, defaults to repository's default branch)

