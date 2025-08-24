# Test script for envm project (PowerShell)
param(
    [switch]$Bench,
    [switch]$Integration,
    [switch]$NoCleanup,
    [int]$CoverageThreshold = 70,
    [switch]$Help
)

# Configuration - Following CONFIGURATION.md standards
$TestDir = "$env:USERPROFILE\.envm-test-$(Get-Random)"
$script:COVERAGE_THRESHOLD = $CoverageThreshold

# Functions
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Setup-TestEnv {
    Write-Status "Setting up test environment..."
    
    # Create test directories
    New-Item -ItemType Directory -Path "$TestDir\downloads\go" -Force | Out-Null
    New-Item -ItemType Directory -Path "$TestDir\downloads\java" -Force | Out-Null
    
    # Set environment variables following CONFIGURATION.md format
    $env:ENVM_HOME = $TestDir
    $env:ENVM_GO_SYMLINK = "$TestDir\go"
    $env:ENVM_JAVA_SYMLINK = "$TestDir\java"
    
    Write-Status "Test environment created at: $TestDir"
}

function Cleanup-TestEnv {
    Write-Status "Cleaning up test environment..."
    if (Test-Path $TestDir) {
        Remove-Item -Path $TestDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    Remove-Item -Path "Env:\ENVM_HOME" -ErrorAction SilentlyContinue
    Remove-Item -Path "Env:\ENVM_GO_SYMLINK" -ErrorAction SilentlyContinue
    Remove-Item -Path "Env:\ENVM_JAVA_SYMLINK" -ErrorAction SilentlyContinue
}

function Run-UnitTests {
    Write-Status "Running unit tests..."
    
    # Run tests with coverage
    $result = go test -v -race -coverprofile=coverage.out ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Unit tests failed"
        return $false
    }
    
    # Generate coverage report
    go tool cover -html=coverage.out -o coverage.html
    
    # Check coverage threshold
    $coverageOutput = go tool cover -func=coverage.out | Select-String "total"
    if ($coverageOutput) {
        $coverage = [double]($coverageOutput.Line -split '\s+')[-1].Replace('%', '')
        if ($coverage -lt $script:COVERAGE_THRESHOLD) {
            Write-Warning "Coverage $coverage% is below threshold $($script:COVERAGE_THRESHOLD)%"
        } else {
            Write-Status "Coverage: $coverage% (above threshold)"
        }
    }
    
    return $true
}

function Run-IntegrationTests {
    Write-Status "Running integration tests..."
    
    # Run integration tests (if any)
    $result = go test -tags=integration -v ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Warning "Integration tests failed or not found"
    }
}

function Run-Benchmarks {
    Write-Status "Running benchmarks..."
    
    $result = go test -bench=. -benchmem ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Warning "Benchmarks failed or not found"
    }
}

function Check-Dependencies {
    Write-Status "Checking dependencies..."
    
    # Check Go version
    go version
    
    # Verify modules
    go mod verify
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Module verification failed"
        return $false
    }
    
    # Download dependencies
    go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to download dependencies"
        return $false
    }
    
    return $true
}

function Show-Help {
    Write-Host @"
Usage: .\test.ps1 [options]

Options:
  -Bench                 Run benchmarks
  -Integration          Run integration tests
  -NoCleanup           Don't cleanup test environment
  -CoverageThreshold N Set coverage threshold (default: 70)
  -Help                Show this help message

Examples:
  .\test.ps1                           # Run basic unit tests
  .\test.ps1 -Bench                   # Run tests with benchmarks
  .\test.ps1 -Integration             # Run tests with integration tests
  .\test.ps1 -CoverageThreshold 80    # Set coverage threshold to 80%
  .\test.ps1 -NoCleanup               # Keep test environment after completion
"@
}

# Main execution
function Main {
    if ($Help) {
        Show-Help
        return
    }
    
    # Setup cleanup if needed
    if (-not $NoCleanup) {
        Register-EngineEvent PowerShell.Exiting -Action { Cleanup-TestEnv }
        trap { Cleanup-TestEnv; break }
    }
    
    try {
        # Execute test pipeline
        if (-not (Check-Dependencies)) {
            exit 1
        }
        
        Setup-TestEnv
        
        if (-not (Run-UnitTests)) {
            exit 1
        }
        
        if ($Integration) {
            Run-IntegrationTests
        }
        
        if ($Bench) {
            Run-Benchmarks
        }
        
        Write-Status "All tests completed successfully!"
        
    } catch {
        Write-Error "Test execution failed: $_"
        exit 1
    } finally {
        if (-not $NoCleanup) {
            Cleanup-TestEnv
        }
    }
}

# Run main function
Main