# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ENVM is a Go-based environment version manager for Go and Java, similar to nvm for Node.js. It allows users to install, switch between, and manage multiple versions of Go and Java runtimes.

## Development Commands

### Building
- `go build -o bin/envm .` - Build the main executable
- `make build` - Build using Makefile (outputs to `bin/envm`)

### Testing
- `make test` - Run all tests
- `make test-verbose` - Run tests with verbose output  
- `make test-cover` - Run tests with coverage report (generates `coverage.html`)
- `make test-race` - Run tests with race detector
- `make test-short` - Run short tests only (skip long-running tests)
- `make test-unit` - Run unit tests only (exclude network tests)
- `make test-windows` - Run tests using PowerShell script (Windows)
- `make test-unix-direct` - Run tests with Unix environment variables
- Individual package tests: `make test-cmd`, `make test-config`, `make test-arch`, `make test-commands`, `make test-util`

### Development Tools
- `make lint` - Run linters (requires golangci-lint)
- `make deps` - Download and tidy Go modules
- `make clean` - Clean test artifacts and cache
- `make ci-test` - Run tests in CI environment with race detection and coverage

## Architecture

### Core Structure
- **main.go** - Entry point that delegates to `cmd.Execute()`
- **cmd/** - CLI command definitions using urfave/cli framework
  - `root.go` - Main app setup and command routing
  - `base.go` - Command definitions for go/java subcommands
- **internal/config/** - Configuration management and environment validation
- **internal/commands/** - Command implementations split by language
  - `commands-go/` - Go version management commands (ls, lsr, active, install, uninstall)
  - `commands-java/` - Java version management commands (ls, active, uninstall)
  - `common/` - Shared command logic
- **internal/arch/** - System architecture detection (unix/windows)
- **internal/logic/web-*/** - Version collection from remote sources
- **util/** - File operations, downloads, and version handling

### Key Environment Variables
- `ENVM_HOME` - Root directory for envm (required)
- `ENVM_GO_SYMLINK` - Symlink location for active Go version
- `ENVM_JAVA_SYMLINK` - Symlink location for active Java version

### Command Structure
The CLI follows a nested structure:
```
envm
├── arch                    # Show system architecture
├── go
│   ├── ls                  # List installed Go versions
│   ├── lsr [stable|archived]  # List remote Go versions
│   ├── active <version>    # Switch to Go version
│   ├── install <version>   # Install Go version
│   └── uninstall <version> # Uninstall Go version
└── java
    ├── ls                  # List installed Java versions
    ├── active <version>    # Switch to Java version
    └── uninstall <version> # Uninstall Java version
```

### Testing Strategy
- Unit tests for individual packages
- Integration tests with proper environment setup
- Platform-specific test scripts (Windows PowerShell, Unix bash)
- Coverage reporting and race condition detection
- Network-dependent tests can be excluded with `-tags=unit`

### Dependencies
- **urfave/cli** - CLI framework
- **mholt/archiver/v3** - Archive extraction
- **PuerkitoBio/goquery** - HTML parsing for version collection
- **blang/semver/v4** - Semantic version handling
- **smartystreets/goconvey** - Testing framework

## Development Notes

- Configuration validation happens before command execution via `app.Before`
- Each language command group has its own environment validation
- Cross-platform support with OS-specific implementations in `internal/arch/`
- Network operations for fetching available versions are in `internal/logic/web-*/`
- File operations and version management utilities are centralized in `util/`