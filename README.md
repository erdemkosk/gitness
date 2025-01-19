# Gitness

Gitness is your repository's fitness tracker! Measure your project's health by discovering its bus factor and contributor dynamics.

## Features

- Calculate repository bus factor
- Analyze contributor statistics
- Support for multiple VCS providers (GitHub, Bitbucket)
- Multiple output formats (Console, JSON, Markdown)
- CI/CD pipeline integration
- Docker support

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
## Usage

### Command Line
```bash
gitness https://github.com/user/repo
```
### With output format
```bash
gitness --output json https://github.com/user/repo
```
### Using environment variables
```bash
export REPOSITORY_URL="https://github.com/user/repo"
export GITHUB_TOKEN="your_token"
gitness
```

### Environment Variables

- `REPOSITORY_URL`: Target repository URL
- `GITHUB_TOKEN`: GitHub personal access token
- `BITBUCKET_USERNAME`: Bitbucket username
- `BITBUCKET_PASSWORD`: Bitbucket app password
- `OUTPUT_FORMAT`: Output format (console, json, markdown)

### CI/CD Integration

#### GitHub Actions

```yaml
name: Gitness Analysis
on: [push]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Gitness
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          docker run -e GITHUB_TOKEN=$GITHUB_TOKEN \
                    -e REPOSITORY_URL=$GITHUB_REPOSITORY \
                    gitness --output json
```

## Output Formats

### Console (default)
```
Repo: user/repo
Bus Factor: 3

Contributors:
------------------
John Doe: 150 commits (45.5%)
Jane Smith: 100 commits (30.3%)
Bob Johnson: 80 commits (24.2%)
```

### JSON
```json
{
  "owner": "user",
  "repo": "repo",
  "busFactor": 3,
  "totalCommits": 330,
  "contributors": [
    {
      "name": "John Doe",
      "commits": 150,
      "percentage": 45.5
    }
  ]
}
```

### Markdown
```markdown
# Repository Analysis: user/repo

## Bus Factor: 3

Total Commits: 330

## Contributors

| Name | Commits | Percentage |
|------|---------|------------|
| John Doe | 150 | 45.5% |
```

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
