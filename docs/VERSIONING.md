# Versioning and Release Management

## Overview

Knetz follows [Semantic Versioning 2.0.0](https://semver.org/) and uses automated version bumping based on [Conventional Commits](https://www.conventionalcommits.org/).

## Semantic Versioning

Given a version number `MAJOR.MINOR.PATCH`, we increment:

- **MAJOR** version when making incompatible API changes
- **MINOR** version when adding functionality in a backwards compatible manner
- **PATCH** version when making backwards compatible bug fixes

## Conventional Commits

All commit messages must follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: A new feature (triggers MINOR bump)
- **fix**: A bug fix (triggers PATCH bump)
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement (triggers PATCH bump)
- **test**: Adding or updating tests
- **build**: Changes to build system or dependencies
- **ci**: Changes to CI/CD configuration
- **chore**: Other changes that don't modify src or test files

### Breaking Changes

For **MAJOR** version bumps, use one of these formats:

**Option 1: Add `!` after type**
```
feat!: change dependency specification format

BREAKING CHANGE: The dependency spec now uses a new schema
```

**Option 2: Use footer**
```
feat(api): update authentication method

BREAKING CHANGE: Authentication now requires OAuth tokens instead of API keys
```

### Examples

```bash
# Minor version bump (new feature)
feat(cluster): add AWS EKS auto-discovery

# Patch version bump (bug fix)
fix(deps): resolve version comparison edge case

# Patch version bump (performance)
perf(scan): optimize cluster scanning with parallel execution

# Major version bump (breaking change)
feat!: redesign configuration file format

BREAKING CHANGE: Configuration format changed from JSON to YAML
```

## Automated Release Process

### Automatic Version Bump (on main branch push)

When commits are pushed to the `main` branch:

1. **Analyze Commits**: Workflow analyzes all commits since the last tag
2. **Determine Bump Type**:
   - `BREAKING CHANGE` or `!` → MAJOR bump
   - `feat:` → MINOR bump
   - `fix:`, `perf:`, etc. → PATCH bump
3. **Generate Changelog**: Creates/updates `CHANGELOG.md` with categorized changes
4. **Update Version**: Updates `VERSION` file and README
5. **Create Tag**: Creates annotated Git tag (e.g., `v1.2.3`)
6. **Push Changes**: Pushes commit and tag to repository
7. **Create Release**: Creates GitHub release with generated notes
8. **Trigger Build**: Tag push triggers GoReleaser to build binaries

### Manual Version Bump

You can manually trigger a version bump from the GitHub Actions UI:

1. Go to **Actions** → **Version Bump and Release**
2. Click **Run workflow**
3. Select bump type: `auto`, `major`, `minor`, or `patch`
4. Click **Run workflow**

### Local Version Bump (for testing)

```bash
# Bump patch version (0.1.0 → 0.1.1)
./scripts/version-bump.sh patch

# Bump minor version (0.1.0 → 0.2.0)
./scripts/version-bump.sh minor

# Bump major version (0.1.0 → 1.0.0)
./scripts/version-bump.sh major
```

## Changelog Generation

The changelog is automatically generated using `git-chglog` with the following categories:

### Features
All commits with `feat:` type

### Bug Fixes
All commits with `fix:` type

### Performance Improvements
All commits with `perf:` type

### Code Refactoring
All commits with `refactor:` type

### Documentation
All commits with `docs:` type

### Others
All other commit types

### Breaking Changes
Any commit with `BREAKING CHANGE:` in footer or `!` after type

## Pull Request Requirements

### PR Title Format

All PR titles must follow Conventional Commits format:

✅ **Valid PR Titles:**
```
feat(cluster): add GCP GKE integration
fix(deps): handle circular dependency detection
docs(readme): update installation instructions
perf(scan): improve scanning performance by 50%
```

❌ **Invalid PR Titles:**
```
Add GCP support
Fixed bug
Update docs
Improvements
```

### PR Validation

A GitHub Action automatically validates PR titles. PRs with invalid titles will fail the check.

### Creating a PR

1. Use the PR template provided
2. Follow Conventional Commits format in the PR title
3. Mark breaking changes with `!` if applicable
4. Fill out all sections of the PR template
5. Link related issues

## Release Workflow

### 1. Development

```bash
# Create feature branch
git checkout -b feat/aws-eks-integration

# Make changes and commit with conventional format
git commit -m "feat(cloud): add AWS EKS cluster discovery"

# Push and create PR
git push origin feat/aws-eks-integration
```

### 2. Pull Request

- PR title must follow conventional commits format
- PR description filled out completely
- All checks pass (tests, linting, PR title validation)
- Code review approved

### 3. Merge to Main

```bash
# Merge PR (via GitHub UI or CLI)
# Automated workflow triggers on main branch push
```

### 4. Automatic Release

The workflow automatically:
- Analyzes commits
- Determines version bump
- Generates changelog
- Creates tag
- Publishes release
- Builds binaries (via GoReleaser)

### 5. Release Artifacts

Each release includes:
- Pre-built binaries for Linux, macOS, Windows (amd64, arm64)
- Checksums file
- Release notes with categorized changes
- Full changelog
- Docker images (future)

## Version Files

### VERSION
Contains the current version number without the `v` prefix:
```
0.1.0
```

### CHANGELOG.md
Contains the full project changelog with all versions and changes.

### Git Tags
All versions are tagged in Git with the `v` prefix:
```
v0.1.0
v0.2.0
v1.0.0
```

## Configuration Files

### .chglog/config.yml
Configuration for `git-chglog` changelog generator

### .chglog/CHANGELOG.tpl.md
Template for changelog generation

### .github/workflows/version-bump.yml
Automated version bump workflow

### .github/workflows/conventional-commits.yml
PR title validation workflow

### .goreleaser.yml
Release configuration with enhanced changelog grouping

## Best Practices

### Commit Messages

✅ **Do:**
- Use present tense ("add feature" not "added feature")
- Use imperative mood ("move cursor to..." not "moves cursor to...")
- Keep subject line under 50 characters
- Capitalize first letter of subject
- Don't end subject with a period
- Use body to explain what and why (not how)

❌ **Don't:**
- Use past tense
- Use vague messages like "fix bug" or "update code"
- Include multiple types in one commit
- Forget to add scope when it adds clarity

### Scopes

Common scopes in Knetz:
- `cluster`: Cluster connection and management
- `scan`: Service discovery and scanning
- `deps`: Dependency management
- `graph`: Dependency graph and analysis
- `cloud`: Cloud provider integration (AWS, GCP, Azure)
- `config`: Configuration management
- `storage`: Database and storage
- `cli`: CLI commands and interface
- `ci`: CI/CD pipelines
- `docs`: Documentation

### Breaking Changes

Only introduce breaking changes when absolutely necessary:
- Changing API interfaces
- Removing features
- Changing configuration format
- Changing CLI command structure
- Changing data storage format

Always provide migration guides for breaking changes.

## Troubleshooting

### Workflow Fails to Create Tag

**Issue:** Tag already exists

**Solution:**
```bash
# Delete local and remote tag
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0

# Workflow will recreate it
```

### PR Title Validation Fails

**Issue:** PR title doesn't follow conventional commits

**Solution:**
Edit PR title to match format: `type(scope): subject`

### Version Not Bumping

**Issue:** No conventional commits found since last tag

**Solution:**
- Ensure at least one commit follows conventional commits format
- Check that commits are of types that trigger bumps (feat, fix, perf, etc.)
- Manually trigger workflow with specific bump type

### Changelog Missing Commits

**Issue:** Some commits not in changelog

**Solution:**
- Commits with types in `.goreleaser.yml` filters are excluded
- Ensure commit follows conventional commits format
- Check that commit type is not in exclude list

## References

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [git-chglog](https://github.com/git-chglog/git-chglog)
- [GoReleaser](https://goreleaser.com/)
- [Keep a Changelog](https://keepachangelog.com/)

