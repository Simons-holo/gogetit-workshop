# Contributing to gogetit-workshop

Thank you for your interest in contributing to this educational project. This guide will help you through the contribution process.

## Contribution Workflow

### 1. Fork the Repository

Click the "Fork" button at the top right of the repository page to create your own copy.

### 2. Clone Your Fork

```bash
git clone https://github.com/YOUR_USERNAME/gogetit-workshop.git
cd gogetit-workshop
```

### 3. Create a Branch

```bash
git checkout -b fix/issue-number-short-description
```

Branch naming conventions:
- `fix/` - Bug fixes
- `feat/` - New features
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### 4. Make Your Changes

Write clean, well-documented code. Follow the code style guidelines below.

### 5. Test Your Changes

```bash
# Run all tests
go test ./...

# Run linter
golangci-lint run

# Format code
gofmt -s -w .
```

### 6. Commit Your Changes

Use conventional commit messages. See the Conventional Commits section below.

```bash
git add .
git commit -m "fix: resolve race condition in downloader"
```

### 7. Push to Your Fork

```bash
git push origin fix/issue-number-short-description
```

### 8. Open a Pull Request

Go to the original repository and click "New Pull Request". Select your branch and fill out the PR template.

## Conventional Commits

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only changes |
| `style` | Changes that do not affect code meaning (formatting, etc.) |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or modifying tests |
| `chore` | Changes to build process or auxiliary tools |

### Examples

```
fix: resolve race condition in concurrent downloader

feat(download): add retry logic for failed downloads

docs: update installation instructions in README

test(scraper): add tests for HTML parsing edge cases
```

### Rules

- Use lowercase for type and scope
- Description should be at least 10 characters
- Use imperative mood ("add" not "added" or "adds")
- No period at the end of the description

## Code Style

### Formatting

All code must be formatted with `gofmt`:

```bash
gofmt -s -w .
```

### Linting

We use `golangci-lint` for static analysis:

```bash
# Install golangci-lint (if not installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

### Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go) principles
- Exported functions must have documentation comments
- Handle all errors explicitly (no ignored errors)
- Use meaningful variable names
- Keep functions focused and small
- Write table-driven tests

## Testing Requirements

- All new code must have tests
- Bug fixes must include a test that reproduces the bug
- Run `go test ./...` before submitting
- Test coverage should not decrease

### Writing Good Tests

```go
func TestDownload(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    []byte
        wantErr bool
    }{
        {
            name:    "valid URL",
            url:     "https://example.com/file.txt",
            want:    []byte("content"),
            wantErr: false,
        },
        {
            name:    "invalid URL",
            url:     "not-a-url",
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Download(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !bytes.Equal(got, tt.want) {
                t.Errorf("Download() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Common Mistakes to Avoid

1. **Not reading the issue carefully**: Understand the problem before starting

2. **Skipping tests**: Always write and run tests for your changes

3. **Ignoring the linter**: Fix all linting errors before submitting

4. **Large PRs**: Keep pull requests small and focused

5. **Missing documentation**: Update docs for changed behavior

6. **Wrong branch name**: Follow the branch naming conventions

7. **Bad commit messages**: Follow conventional commit format

8. **Not rebasing on main**: Keep your branch up to date

9. **Unrelated changes**: One PR should address one issue

10. **Not checking CI results**: Ensure all checks pass

## Keeping Your Fork Updated

```bash
# Add upstream remote
git remote add upstream https://github.com/anxkhn/gogetit-workshop.git

# Fetch upstream changes
git fetch upstream

# Rebase your branch on upstream main
git checkout main
git merge upstream/main
git checkout your-branch
git rebase main
```

## How to Get Help

- **Issues**: Comment on the issue you're working on with questions
- **Discussions**: Use GitHub Discussions for general questions
- **Pull Requests**: Open a draft PR early to get feedback

## Pull Request Template

When you open a pull request, fill out the template completely. See `.github/PULL_REQUEST_TEMPLATE.md` for the template structure.

## Code of Conduct

Be respectful and constructive. We're all here to learn and improve.
