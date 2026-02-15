## Description

<!-- Provide a clear description of your changes -->


## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)
- [ ] Test addition or modification

## Checklist

- [ ] I have read the CONTRIBUTING.md document
- [ ] My code follows the code style of this project
- [ ] I have run `gofmt -s -w .` on my code
- [ ] I have run `golangci-lint run` and fixed all issues
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with `go test ./...`
- [ ] I have updated the documentation accordingly

## Commit Message Format

This PR follows conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for test changes
- `chore:` for maintenance tasks

Example: `fix: resolve race condition in concurrent downloader`

## Golangci-lint

- [ ] I have run `golangci-lint run` locally and all checks pass

## Related Issues

<!-- Link any related issues here using #issue-number -->
