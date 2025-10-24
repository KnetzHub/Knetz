# Contributing to Knetz

Thank you for your interest in contributing to Knetz! This document provides guidelines and instructions for contributing.

## 🤝 Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- Git
- Access to a Kubernetes cluster (for testing)

### Development Setup

1. Fork and clone the repository:
```bash
git clone https://github.com/your-username/knetz.git
cd knetz
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
make build
```

4. Run tests:
```bash
make test
```

## 📝 How to Contribute

### Reporting Bugs

- Check if the bug has already been reported in [GitHub Issues](https://github.com/knetz-io/knetz/issues)
- Use the bug report template
- Include:
  - Clear description of the issue
  - Steps to reproduce
  - Expected vs actual behavior
  - Environment details (OS, Go version, Kubernetes version)
  - Relevant logs or screenshots

### Suggesting Features

- Open a GitHub Issue with the "enhancement" label
- Clearly describe the feature and its use case
- Explain how it aligns with Knetz's goals

### Pull Requests

1. Create a new branch from `main`:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes following our coding standards

3. Write or update tests

4. Ensure all tests pass:
```bash
make test
make lint
```

5. Commit your changes with clear messages:
```bash
git commit -m "feat: add support for X"
```

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

6. Push and create a pull request

7. Describe your changes in the PR description

## 🎯 Coding Standards

### Go Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and concise

### Testing

- Write unit tests for new functionality
- Aim for 80%+ code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

### Documentation

- Update README.md for user-facing changes
- Add/update inline documentation
- Include examples for new features

## 🏗️ Project Structure

```
knetz/
├── cmd/knetz/           # CLI entry point and commands
│   ├── commands/        # Cobra command implementations
│   └── main.go
├── pkg/                 # Public packages
│   ├── cluster/         # Cluster connection management
│   ├── config/          # Configuration handling
│   ├── discovery/       # Service discovery
│   ├── storage/         # Data persistence
│   └── ...
├── internal/            # Private packages
│   ├── models/          # Data models
│   └── utils/           # Utility functions
└── docs/                # Documentation

```

## 🧪 Testing Guidelines

### Unit Tests

```go
func TestVersionCompare(t *testing.T) {
    tests := []struct {
        name     string
        v1       string
        v2       string
        expected int
    }{
        {"equal versions", "1.0.0", "1.0.0", 0},
        {"v1 greater", "1.1.0", "1.0.0", 1},
        {"v2 greater", "1.0.0", "1.1.0", -1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CompareVersionStrings(tt.v1, tt.v2)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Tests

- Use kind or k3s for Kubernetes testing
- Clean up resources after tests
- Test against multiple Kubernetes versions

## 🔄 Release Process

Releases are handled by maintainers:

1. Update version in appropriate files
2. Update CHANGELOG.md
3. Create and push a git tag
4. GitHub Actions will build and publish the release

## 📞 Getting Help

- GitHub Discussions for questions
- Slack channel (coming soon)
- Tag maintainers in issues/PRs if needed

## 🎖️ Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project README

Thank you for contributing to Knetz! 🚀

