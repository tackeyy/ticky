# Contributing to ticky

Thank you for your interest in contributing to ticky! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Ways to Contribute](#ways-to-contribute)
- [Before You Start](#before-you-start)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Submitting Changes](#submitting-changes)
- [Code Review Process](#code-review-process)
- [Community Guidelines](#community-guidelines)
- [Getting Help](#getting-help)

## Ways to Contribute

### You can contribute by:

- 🐛 **Reporting bugs** - Found an issue? Let us know!
- 💡 **Suggesting features** - Have an idea? We'd love to hear it
- 📝 **Improving documentation** - Help make our docs clearer
- 🔧 **Submitting bug fixes** - Fix issues and help improve stability
- ✨ **Adding new features** - Expand ticky's capabilities (discuss first!)

## Before You Start

1. **Check existing issues/PRs** to avoid duplication
2. **For new features**, open an issue first to discuss the proposal
3. **Read our [Testing Guide](docs/TESTING.md)** to understand our testing approach
4. **Ensure you understand our [Code of Conduct](CODE_OF_CONDUCT.md)**

## Development Setup

### Prerequisites

- Go 1.25.2+
- TickTick OAuth credentials ([How to get credentials](https://developer.ticktick.com/manage))

### Setup Steps

```bash
# 1. Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/ticky.git
cd ticky

# 2. Set up environment variables
export TICKTICK_CLIENT_ID=your_client_id
export TICKTICK_CLIENT_SECRET=your_client_secret

# 3. Run tests to verify setup
go test -v ./...

# 4. Build the project
go build -o ticky .

# 5. Test the CLI locally
./ticky --help
```

### TickTick App Setup

1. Go to [TickTick Developer Portal](https://developer.ticktick.com/manage) and click **+Create App**
2. Name your app (e.g., `ticky-dev`)
3. Set the Redirect URL to `http://localhost:18080/callback`
4. Enable scopes: `tasks:read`, `tasks:write`

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format your code (already enforced by Go tooling)
- Use descriptive variable names (`projectID` not `id`)
- Add comments for exported functions and types
- Keep functions small and focused (single responsibility)
- Handle errors explicitly; do not ignore them

### Code Organization

- CLI commands go in `cmd/`
- TickTick API client and business logic go in `internal/ticktick/`
- Follow existing patterns in the codebase

### Commit Message Convention

Format: `<type>: <subject>`

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `test:` Test additions/changes
- `docs:` Documentation changes
- `refactor:` Code refactoring (no functional changes)
- `chore:` Maintenance tasks (dependencies, tooling)

**Examples:**
```
feat: add support for recurring tasks
fix: correct date parsing for relative format
test: add validation tests for priority parser
docs: update README with new --tags flag
refactor: extract HTTP helper to separate function
chore: update cobra to latest version
```

## Testing Requirements

**All code contributions MUST include tests.**

### Test Types

1. **Unit Tests** - Test individual functions in isolation (e.g., `ParsePriority`, `ParseDate`)
2. **HTTP Mock Tests** - Test API client methods using `httptest.NewServer` (e.g., `GetProjects`, `CreateTask`)
3. **Authentication Tests** - Test OAuth token exchange and refresh flows
4. **File I/O Tests** - Test token storage and configuration using `t.TempDir()`

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v -run TestParsePriority ./internal/ticktick/
```

### Test Writing Guidelines

- Follow **Arrange/Act/Assert** pattern
- Use **table-driven tests** with subtests (`t.Run`)
- Use descriptive test names: `TestParsePriority_ValidInputs`, `TestExchangeToken_HTTPError`
- Mock HTTP endpoints with `net/http/httptest`
- Use `t.TempDir()` and `t.Setenv()` for file system and environment isolation
- See **[docs/TESTING.md](docs/TESTING.md)** for the comprehensive testing guide

### Test Coverage Expectations

- **New features**: Tests required for all new code
- **Bug fixes**: Add regression test reproducing the bug
- **Refactoring**: Maintain or improve existing coverage
- **Target**: 80%+ for testable code

## Submitting Changes

### Pull Request Process

#### 1. Create a branch

```bash
git checkout -b feat/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

#### 2. Make your changes

- Write code
- Add tests
- Update documentation if needed

#### 3. Ensure quality

```bash
go test -v ./...     # All tests must pass
go build -o ticky .  # Build must succeed
```

#### 4. Commit your changes

```bash
git add .
git commit -m "feat: add your feature description"
```

#### 5. Push and create PR

```bash
git push origin feat/your-feature-name
# Then create PR via GitHub UI
```

#### 6. Fill out PR template

- Describe what changed and why
- Link related issues with `Closes #123`
- Provide testing evidence
- Check all applicable boxes in the template

### PR Requirements Checklist

Before submitting, ensure:

- All tests pass (`go test -v ./...`)
- Build succeeds (`go build -o ticky .`)
- Code follows project style (`gofmt`)
- Commit messages follow convention
- Tests added for new functionality
- Documentation updated (if applicable)
- PR template fully completed

### What to Expect

- **Initial review** within 2-3 business days
- **Feedback** and requested changes from maintainers
- **Approval and merge** once all requirements are met

## Code Review Process

### For Contributors

- **Be responsive** to feedback and questions
- **Ask for clarification** if feedback is unclear
- **Push updates** to the same branch (PR will auto-update)
- **Be patient and respectful** throughout the process

### Review Criteria

Reviewers will check:

- **Functionality** - Does it work as intended?
- **Tests** - Are they comprehensive and passing?
- **Code Quality** - Is it readable and maintainable?
- **Documentation** - Is it clear and up-to-date?
- **Performance** - Are there any obvious performance issues?
- **Security** - Are there any potential vulnerabilities?

## Community Guidelines

- Be respectful and welcoming to all contributors
- Follow our [Code of Conduct](CODE_OF_CONDUCT.md)
- Provide constructive feedback
- Assume good intentions
- Help others learn and grow

## Getting Help

- **Bug Reports** - Open an [Issue](https://github.com/tackeyy/ticky/issues/new?template=bug_report.yml)
- **Feature Requests** - Open an [Issue](https://github.com/tackeyy/ticky/issues/new?template=feature_request.yml)
- **General Questions** - Open an [Issue](https://github.com/tackeyy/ticky/issues)

## Recognition

All contributors are recognized in:

- GitHub Contributors page
- Release notes (for significant contributions)

---

Thank you for contributing to ticky! Your efforts help make this tool better for everyone.
