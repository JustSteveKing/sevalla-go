# Sevalla Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/juststeveking/sevalla-go.svg)](https://pkg.go.dev/github.com/juststeveking/sevalla-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/juststeveking/sevalla-go)](https://goreportcard.com/report/github.com/juststeveking/sevalla-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go client library for the Sevalla API, providing type-safe access to all Sevalla services including applications, databases, static sites, deployments, and CI/CD pipelines.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Core Concepts](#core-concepts)
- [Usage Guide](#usage-guide)
  - [Applications](#applications)
  - [Databases](#databases)
  - [Static Sites](#static-sites)
  - [Deployments](#deployments)
  - [Pipelines](#pipelines)
- [Best Practices](#best-practices)
- [Error Handling](#error-handling)
- [Advanced Topics](#advanced-topics)
- [Testing](#testing)
- [Examples](#examples)
- [Contributing](#contributing)

## Features

- ✅ **Type-safe API** - Strongly typed requests and responses with zero reflection
- ✅ **Full API coverage** - Support for all Sevalla API endpoints
- ✅ **Production-ready** - 79.7% test coverage with comprehensive test suite
- ✅ **Error handling** - Rich error types with helper functions for common scenarios
- ✅ **Pagination support** - Built-in pagination handling for list operations
- ✅ **Context support** - Full context.Context support for timeouts and cancellation
- ✅ **Clean API design** - Intuitive service-based architecture following Go best practices
- ✅ **Zero dependencies** - Only standard library and minimal external dependencies
- ✅ **Retry logic ready** - Structured for easy integration with retry mechanisms
- ✅ **Mock-friendly** - Interface-based design for easy testing

## Installation

```bash
go get github.com/juststeveking/sevalla-go
```

Requires Go 1.19 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
    // Create a new client with your API key
    client := sevalla.NewClient(
        sevalla.WithAPIKey("your-api-key"),
    )
    
    ctx := context.Background()
    
    // List all applications
    apps, resp, err := client.Applications.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d applications (Status: %d)\n", len(apps), resp.StatusCode)
    
    for _, app := range apps {
        fmt.Printf("  - %s (State: %s, URL: %s)\n", app.Name, app.State, app.URL)
    }
}
```

## Authentication

The SDK requires a Sevalla API key for authentication. You can obtain your API key from the [Sevalla Dashboard](https://app.sevalla.com/company/apiKeys).

**Best Practice:** Store your API key securely using environment variables:

```go
package main

import (
    "os"
    sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
    apiKey := os.Getenv("SEVALLA_API_KEY")
    if apiKey == "" {
        log.Fatal("SEVALLA_API_KEY environment variable is required")
    }
    
    client := sevalla.NewClient(
        sevalla.WithAPIKey(apiKey),
    )
    
    // Use the client...
}
```

**Security Note:** Never commit API keys to version control. Use environment variables, secret managers, or configuration management tools.

## Core Concepts

### Client Architecture

The SDK is organized into service-based clients, each handling a specific resource type:

```go
client := sevalla.NewClient(sevalla.WithAPIKey("key"))

// Available services:
client.Applications  // Manage applications
client.Databases     // Manage databases
client.StaticSites   // Manage static sites
client.Deployments   // Monitor deployments
client.Pipelines     // Manage CI/CD pipelines
```

### Context Usage

All API methods accept a `context.Context` as their first parameter, following Go best practices:

```go
// Basic context
ctx := context.Background()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
defer cancel()
```

### Response Structure

All API methods return three values:

1. **Result** - The decoded response data (e.g., `*Application`, `[]*Database`)
2. **Response** - HTTP response metadata including pagination info
3. **Error** - Any error that occurred during the request

```go
app, resp, err := client.Applications.Get(ctx, "app-123")
if err != nil {
    // Handle error
}

// resp contains:
// - resp.StatusCode (int)
// - resp.NextPage (int)
// - resp.PrevPage (int)
// - resp.FirstPage (int)
// - resp.LastPage (int)
```

## Usage Guide

### Applications

Applications represent your deployed services, web apps, or APIs running on Sevalla infrastructure.

#### Creating an Application

```go
createReq := &sevalla.CreateApplicationRequest{
    Name:          "my-production-app",
    RepositoryURL: "https://github.com/user/my-app",
    Branch:        "main",
    Region:      sevalla.RegionUSCentral,
    Plan:       sevalla.PlanStarter,
    Replicas:      2,
    BuildCommand:  "npm run build",
    StartCommand:  "npm start",
    Port:          3000,
    AutoDeploy:    true,
    CDNEnabled:    true,
    SSLEnabled:    true,
    EnvironmentVars: map[string]string{
        "NODE_ENV":    "production",
        "DATABASE_URL": "postgresql://...",
        "API_KEY":     os.Getenv("API_KEY"),
    },
}

app, resp, err := client.Applications.Create(ctx, createReq)
if err != nil {
    if sevalla.IsBadRequest(err) {
        log.Printf("Invalid request: %v", err)
        return
    }
    log.Fatal(err)
}

fmt.Printf("Created application: %s (ID: %s)\n", app.Name, app.ID)
fmt.Printf("Application URL: %s\n", app.URL)
```

**Best Practice:** Always enable SSL for production applications and use appropriate pod sizes for your workload.

#### Listing Applications

```go
// List all applications
apps, resp, err := client.Applications.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

// With pagination
opts := &sevalla.ListOptions{
    Page:    1,
    PerPage: 25,
    Sort:    "created_at",
    Order:   "desc",
}

apps, resp, err := client.Applications.List(ctx, opts)
```

#### Getting a Single Application

```go
app, resp, err := client.Applications.Get(ctx, "app-123")
if err != nil {
    if sevalla.IsNotFound(err) {
        log.Println("Application not found")
        return
    }
    log.Fatal(err)
}

fmt.Printf("Application: %s\n", app.Name)
fmt.Printf("State: %s\n", app.State)
fmt.Printf("Replicas: %d\n", app.Replicas)
fmt.Printf("Plan: %s\n", app.Plan)
```

#### Updating an Application

```go
updateReq := &sevalla.UpdateApplicationRequest{
    Name:     sevalla.String("my-app-renamed"),
    Replicas: sevalla.Int(3),
    Plan:  &sevalla.PlanPro,
}

app, resp, err := client.Applications.Update(ctx, "app-123", updateReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated to %d replicas\n", app.Replicas)
```

**Note:** Use helper functions `sevalla.String()`, `sevalla.Int()`, and `sevalla.Bool()` for pointer types.

#### Deploying an Application

```go
// Trigger a new deployment
deployment, resp, err := client.Applications.Deploy(ctx, app.ID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deployment started: %s\n", deployment.ID)
fmt.Printf("State: %s\n", deployment.State)

// Monitor deployment status
ticker := time.NewTicker(5 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        dep, _, err := client.Deployments.Get(ctx, deployment.ID)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("Deployment state: %s\n", dep.State)
        
        if dep.State == sevalla.StatusSuccess {
            fmt.Println("Deployment completed successfully!")
            return
        }
        
        if dep.State == sevalla.StatusFailed {
            fmt.Printf("Deployment failed: %s\n", dep.ErrorMessage)
            return
        }
    case <-ctx.Done():
        fmt.Println("Deployment monitoring cancelled")
        return
    }
}
```

#### Scaling an Application

```go
scaleReq := &sevalla.ScaleApplicationRequest{
    Replicas: 5,
}

// Scale to 5 replicas
app, resp, err := client.Applications.Scale(ctx, "app-123", scaleReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Scaled to %d replicas\n", app.Replicas)
```

**Best Practice:** Scale horizontally (more replicas) rather than vertically (larger pods) for better availability.

#### Managing Application Lifecycle

```go
// Stop an application
_, err := client.Applications.Stop(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}

// Start a stopped application
_, err = client.Applications.Start(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}

// Restart an application (useful for applying config changes)
_, err = client.Applications.Restart(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}

// Delete an application (careful - this is permanent!)
_, err = client.Applications.Delete(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}
```

**Warning:** Delete operations are permanent and cannot be undone.

#### Environment Variables

```go
// Set multiple environment variables
vars := map[string]string{
    "DATABASE_URL": "postgresql://new-host/db",
    "API_KEY":      "new-secret-key",
    "DEBUG":        "false",
}

_, err := client.Applications.SetEnvironmentVariables(ctx, "app-123", vars)
if err != nil {
    log.Fatal(err)
}

// Get current environment variables
currentVars, resp, err := client.Applications.GetEnvironmentVariables(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}

for key, value := range currentVars {
    // Note: Sensitive values may be masked
    fmt.Printf("%s: %s\n", key, value)
}
```

**Security:** Environment variables are stored encrypted. Never log sensitive values.

#### Viewing Logs

```go
// Get latest 100 lines
logs, resp, err := client.Applications.GetLogs(ctx, "app-123", 100)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Recent logs:")
fmt.Println(logs)

// Get all available logs
logs, resp, err = client.Applications.GetLogs(ctx, "app-123", 0)
```

#### Managing Custom Domains

```go
// Add a custom domain
_, err := client.Applications.AddCustomDomain(ctx, "app-123", "www.example.com")
if err != nil {
    log.Fatal(err)
}

// Remove a custom domain
_, err = client.Applications.RemoveCustomDomain(ctx, "app-123", "www.example.com")
if err != nil {
    log.Fatal(err)
}
```

#### CDN Settings

```go
// Enable CDN
_, err := client.Applications.UpdateCDNSettings(ctx, "app-123", true)
if err != nil {
    log.Fatal(err)
}

// Disable CDN
_, err = client.Applications.UpdateCDNSettings(ctx, "app-123", false)
```

#### Usage Metrics

```go
// Get usage for the last 7 days
usage, resp, err := client.Applications.GetUsage(ctx, "app-123", "7d")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("CPU Usage: %.2f%%\n", usage.CPUUsage)
fmt.Printf("Memory Usage: %.2f MB\n", usage.MemoryUsage)
fmt.Printf("Bandwidth: %d bytes\n", usage.BandwidthUsed)
fmt.Printf("Requests: %d\n", usage.RequestCount)
```

#### Rollback to Previous Deployment

```go
// List recent deployments
deployments, _, err := client.Applications.ListDeployments(ctx, "app-123", &sevalla.ListOptions{
    PerPage: 10,
})
if err != nil {
    log.Fatal(err)
}

// Find the last successful deployment before current
var previousDeployment *sevalla.Deployment
for _, dep := range deployments {
    if dep.State == sevalla.StatusSuccess && dep.ID != app.LastDeploymentID {
        previousDeployment = dep
        break
    }
}

if previousDeployment != nil {
    // Rollback
    newDeployment, _, err := client.Applications.Rollback(ctx, "app-123", previousDeployment.ID)
    if err != nil {
        log.Fatal(err)
    }
    
        fmt.Printf("Rolled back to deployment %s\n", previousDeployment.ID)
}
```

### Databases

Manage PostgreSQL, MySQL, MongoDB, and Redis databases with full lifecycle control.

#### Creating a Database

```go
createDBReq := &sevalla.CreateDatabaseRequest{
    Name:       "production-db",
    Type:       sevalla.EnginePostgreSQL,
    Version:    "15",  // Latest stable version
    Region:   sevalla.RegionUSCentral,
    Size:       "medium",  // small, medium, large, xlarge
    Storage:    50,        // GB
    Backups:    true,      // Enable automatic backups
    SSLEnabled: true,      // Enforce SSL connections
}

db, resp, err := client.Databases.Create(ctx, createDBReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Database created: %s\n", db.Name)
fmt.Printf("Internal URL: %s\n", db.InternalURL)
fmt.Printf("Type: %s %s\n", db.Type, db.Version)
```

**Best Practice:** Always enable SSL and backups for production databases.

#### Getting Database Credentials

```go
db, resp, err := client.Databases.GetCredentials(ctx, "db-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Username: %s\n", db.Username)
fmt.Printf("Password: %s\n", db.Password)  // Handle securely!
fmt.Printf("Connection URL: %s\n", db.InternalURL)

// Use in your application
connectionString := fmt.Sprintf(
    "postgresql://%s:%s@%s",
    db.Username,
    db.Password,
    strings.TrimPrefix(db.InternalURL, "postgresql://"),
)
```

**Security:** Never log credentials. Store them securely in environment variables or secret management systems.

#### Resetting Database Password

```go
db, resp, err := client.Databases.ResetPassword(ctx, "db-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("New password: %s\n", db.Password)
// Update your application's environment variables immediately
```

#### Managing Backups

```go
// List all backups
backups, resp, err := client.Databases.ListBackups(ctx, "db-123", nil)
if err != nil {
    log.Fatal(err)
}

for _, backup := range backups {
    fmt.Printf("Backup: %s (Created: %s, Size: %d MB, Status: %s)\n",
        backup.ID,
        backup.CreatedAt.Format(time.RFC3339),
        backup.Size/1024/1024,
        backup.Status,
    )
}

// Create a manual backup
backupReq := &sevalla.CreateBackupRequest{
    Type:        "manual",
    Description: "Pre-migration backup",
}

backup, resp, err := client.Databases.CreateBackup(ctx, "db-123", backupReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Backup created: %s\n", backup.ID)

// Restore from backup
restoreReq := &sevalla.RestoreBackupRequest{
    BackupID: backup.ID,
}

_, err = client.Databases.RestoreFromBackup(ctx, "db-123", restoreReq)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Restore initiated. Database will be unavailable during restore.")
```

**Warning:** Restoring a backup will overwrite the current database. Always create a fresh backup before restoring.

#### Monitoring Database Usage

```go
// Get usage for the last 24 hours
usage, resp, err := client.Databases.GetUsage(ctx, "db-123", "24h")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Database Usage:\n")
fmt.Printf("  CPU: %.2f%%\n", usage.CPUUsage)
fmt.Printf("  Memory: %.2f MB\n", usage.MemoryUsage)
fmt.Printf("  Storage: %.2f GB\n", usage.StorageUsage/1024)
fmt.Printf("  Bandwidth: %d MB\n", usage.BandwidthUsed/1024/1024)

// Alert if usage is high
if usage.CPUUsage > 80 {
    log.Println("Warning: High CPU usage detected!")
    // Consider upgrading database size
}
```

#### Public Access Control

```go
// Enable public access (use with caution)
db, resp, err := client.Databases.EnablePublicAccess(ctx, "db-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Public URL: %s\n", db.PublicURL)
fmt.Println("Warning: Database is now accessible from the internet!")

// Disable public access (recommended)
db, resp, err = client.Databases.DisablePublicAccess(ctx, "db-123")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Public access disabled. Database is now only accessible internally.")
```

**Security Best Practice:** Only enable public access for databases that absolutely require it. Use IP whitelisting and strong passwords.

#### Updating Database Configuration

```go
updateReq := &sevalla.UpdateDatabaseRequest{
    Name:    sevalla.String("renamed-db"),
    Size:    sevalla.String("large"),
    Storage: sevalla.Int(100),  // Increase storage to 100GB
    Backups: sevalla.Bool(true),
}

db, resp, err := client.Databases.Update(ctx, "db-123", updateReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Database updated: Size=%s, Storage=%dGB\n", db.Size, db.Storage)
```

**Note:** Some updates may require a restart. Plan accordingly for production databases.

### Static Sites

Deploy and manage static websites and SPAs (Single Page Applications).

#### Creating a Static Site

```go
createSiteReq := &sevalla.CreateStaticSiteRequest{
    Name:            "my-portfolio",
    RepositoryURL:   "https://github.com/user/portfolio",
    Branch:          "main",
    BuildCommand:    "npm run build",
    OutputDirectory: "dist",      // or "build", "_site", "public", etc.
    Region:        sevalla.RegionEuropeWest,
    AutoDeploy:      true,        // Deploy on push
    CDNEnabled:      true,        // Enable CDN for better performance
    SSLEnabled:      true,        // Always use HTTPS
    EnvironmentVars: map[string]string{
        "NODE_ENV":         "production",
        "VITE_API_URL":     "https://api.example.com",
        "NEXT_PUBLIC_KEY":  "public-key",
    },
}

site, resp, err := client.StaticSites.Create(ctx, createSiteReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Static site created: %s\n", site.Name)
fmt.Printf("URL: %s\n", site.URL)
```

**Framework Examples:**
- **Next.js**: `BuildCommand: "npm run build"`, `OutputDirectory: "out"` or `.next`
- **Vite**: `BuildCommand: "npm run build"`, `OutputDirectory: "dist"`
- **Create React App**: `BuildCommand: "npm run build"`, `OutputDirectory: "build"`
- **Jekyll**: `BuildCommand: "bundle exec jekyll build"`, `OutputDirectory: "_site"`
- **Hugo**: `BuildCommand: "hugo"`, `OutputDirectory: "public"`

#### Deploying a Static Site

```go
// Trigger a new build and deployment
deployment, resp, err := client.StaticSites.Deploy(ctx, "site-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deployment started: %s\n", deployment.ID)
fmt.Printf("State: %s\n", deployment.State)
```

#### Listing Static Sites

```go
sites, resp, err := client.StaticSites.List(ctx, &sevalla.ListOptions{
    PerPage: 50,
    Sort:    "name",
})
if err != nil {
    log.Fatal(err)
}

for _, site := range sites {
    fmt.Printf("Site: %s - %s (State: %s)\n", site.Name, site.URL, site.State)
}
```

### Deployments

Monitor and manage application and static site deployments.

#### Getting Deployment Status

```go
deployment, resp, err := client.Deployments.Get(ctx, "deploy-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deployment: %s\n", deployment.ID)
fmt.Printf("State: %s\n", deployment.State)
fmt.Printf("Branch: %s\n", deployment.Branch)
fmt.Printf("Commit: %s\n", deployment.CommitSHA)
fmt.Printf("Started: %s\n", deployment.StartedAt.Format(time.RFC3339))

if deployment.CompletedAt != nil {
    fmt.Printf("Completed: %s\n", deployment.CompletedAt.Format(time.RFC3339))
    fmt.Printf("Duration: %d seconds\n", deployment.Duration)
}

if deployment.State == sevalla.StatusFailed {
    fmt.Printf("Error: %s\n", deployment.ErrorMessage)
}
```

#### Viewing Deployment Logs

```go
logs, resp, err := client.Deployments.GetLogs(ctx, "deploy-123")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Build logs:")
fmt.Println(logs)
```

#### Cancelling a Deployment

```go
_, err := client.Deployments.Cancel(ctx, "deploy-123")
if err != nil {
    if sevalla.IsConflict(err) {
        fmt.Println("Deployment already completed or cancelled")
    } else {
        log.Fatal(err)
    }
}

fmt.Println("Deployment cancelled")
```

#### Listing All Deployments

```go
deployments, resp, err := client.Deployments.List(ctx, &sevalla.ListOptions{
    Page:    1,
    PerPage: 20,
    Sort:    "started_at",
    Order:   "desc",
})
if err != nil {
    log.Fatal(err)
}

for _, dep := range deployments {
    fmt.Printf("%s - %s (State: %s, Branch: %s)\n",
        dep.StartedAt.Format("2006-01-02 15:04"),
        dep.ID[:8],
        dep.State,
        dep.Branch,
    )
}
```

### Pipelines

Create and manage CI/CD pipelines for automated testing and deployment.

#### Creating a Pipeline

```go
pipelineReq := &sevalla.CreatePipelineRequest{
    Name:    "test-and-deploy",
    Enabled: true,
    Trigger: "push",  // or "pull_request", "tag", "manual"
    Branch:  "main",
    Steps: []sevalla.PipelineStep{
        {
            Name:    "install",
            Command: "npm ci",
            Image:   "node:18",
            Timeout: 300,
        },
        {
            Name:      "test",
            Command:   "npm test",
            DependsOn: []string{"install"},
            Timeout:   600,
        },
        {
            Name:      "build",
            Command:   "npm run build",
            DependsOn: []string{"test"},
            Timeout:   600,
        },
        {
            Name:      "deploy",
            Command:   "npm run deploy",
            DependsOn: []string{"build"},
            Timeout:   300,
        },
    },
    Environment: map[string]string{
        "NODE_ENV":  "production",
        "CI":        "true",
    },
}

pipeline, resp, err := client.Pipelines.Create(ctx, pipelineReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Pipeline created: %s (ID: %s)\n", pipeline.Name, pipeline.ID)
```

**Best Practice:** Use step dependencies to ensure proper execution order and fail fast.

#### Triggering a Pipeline Run

```go
run, resp, err := client.Pipelines.Run(ctx, "pipeline-123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Pipeline run started: %s\n", run.ID)
fmt.Printf("State: %s\n", run.State)
```

#### Monitoring Pipeline Runs

```go
// List recent runs
runs, resp, err := client.Pipelines.ListRuns(ctx, "pipeline-123", &sevalla.ListOptions{
    PerPage: 10,
})
if err != nil {
    log.Fatal(err)
}

for _, run := range runs {
    fmt.Printf("Run: %s - State: %s, Branch: %s, Duration: %ds\n",
        run.ID,
        run.State,
        run.Branch,
        run.Duration,
    )
}

// Get specific run details
run, resp, err := client.Pipelines.GetRun(ctx, "pipeline-123", "run-456")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Run %s: %s\n", run.ID, run.State)
for _, step := range run.Steps {
    fmt.Printf("  Step: %s - %s\n", step.Name, step.State)
}
```

#### Viewing Pipeline Logs

```go
logs, resp, err := client.Pipelines.GetRunLogs(ctx, "pipeline-123", "run-456")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Pipeline logs:")
fmt.Println(logs)
```

#### Retrying Failed Runs

```go
newRun, resp, err := client.Pipelines.RetryRun(ctx, "pipeline-123", "run-456")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Retry started: %s\n", newRun.ID)
```

## Best Practices

### 1. Error Handling

Always check and handle errors appropriately:

```go
app, resp, err := client.Applications.Get(ctx, appID)
if err != nil {
    // Check for specific error types
    switch {
    case sevalla.IsNotFound(err):
        return fmt.Errorf("application %s not found", appID)
    case sevalla.IsUnauthorized(err):
        return fmt.Errorf("invalid API key")
    case sevalla.IsRateLimited(err):
        // Implement exponential backoff
        time.Sleep(time.Minute)
        return retry()
    case sevalla.IsServerError(err):
        // Server error - safe to retry
        return retry()
    default:
        return fmt.Errorf("unexpected error: %w", err)
    }
}
```

### 2. Context Management

Use contexts for timeouts and cancellation:

```go
// Set reasonable timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// For long-running operations, use longer timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

// For user-initiated requests, support cancellation
ctx, cancel := context.WithCancel(request.Context())
defer cancel()
```

### 3. Pagination

Handle pagination efficiently:

```go
func listAllApplications(ctx context.Context, client *sevalla.Client) ([]*sevalla.Application, error) {
    var allApps []*sevalla.Application
    
    opts := &sevalla.ListOptions{
        Page:    1,
        PerPage: 100,  // Use maximum page size for fewer requests
    }
    
    for {
        apps, resp, err := client.Applications.List(ctx, opts)
        if err != nil {
            return nil, err
        }
        
        allApps = append(allApps, apps...)
        
        // Check if there are more pages
        if resp.NextPage == 0 {
            break
        }
        
        opts.Page = resp.NextPage
    }
    
    return allApps, nil
}
```

### 4. Resource Cleanup

Always clean up resources in production:

```go
func deployApplication(ctx context.Context, client *sevalla.Client) error {
    // Create temporary test application
    app, _, err := client.Applications.Create(ctx, createReq)
    if err != nil {
        return err
    }
    
    // Ensure cleanup even if errors occur
    defer func() {
        cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if _, err := client.Applications.Delete(cleanupCtx, app.ID); err != nil {
            log.Printf("Failed to cleanup application: %v", err)
        }
    }()
    
    // Use the application...
    return nil
}
```

### 5. Secure Credential Handling

Never expose credentials:

```go
// ❌ Bad: Logging credentials
db, _, err := client.Databases.GetCredentials(ctx, dbID)
log.Printf("Password: %s", db.Password)  // DON'T DO THIS

// ✅ Good: Use credentials securely
db, _, err := client.Databases.GetCredentials(ctx, dbID)
if err != nil {
    return err
}

// Store in environment or secret manager
os.Setenv("DB_PASSWORD", db.Password)  // Or use proper secret management

// Clear from memory when done
db.Password = ""
```

### 6. Retry Logic

Implement exponential backoff for retries:

```go
func withRetry(ctx context.Context, maxRetries int, fn func() error) error {
    backoff := time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Don't retry client errors (4xx) except rate limits
        if sevalla.IsClientError(err) && !sevalla.IsRateLimited(err) {
            return err
        }
        
        // Wait before retrying
        select {
        case <-time.After(backoff):
            backoff *= 2  // Exponential backoff
            if backoff > time.Minute {
                backoff = time.Minute
            }
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return fmt.Errorf("max retries exceeded")
}
```

### 7. Environment-Specific Configuration

Use different configurations for different environments:

```go
func newClient(env string) *sevalla.Client {
    apiKey := os.Getenv("SEVALLA_API_KEY")
    
    opts := []sevalla.ClientOption{
        sevalla.WithAPIKey(apiKey),
    }
    
    // Use custom settings for different environments
    switch env {
    case "development":
        opts = append(opts,
            sevalla.WithHTTPClient(&http.Client{
                Timeout: 60 * time.Second,
            }),
        )
    case "staging":
        opts = append(opts,
            sevalla.WithBaseURL("https://staging-api.sevalla.com/v2"),
        )
    case "production":
        // Use default settings for production
    }
    
    return sevalla.NewClient(opts...)
}
```

### 8. Monitoring and Logging

Log important operations for debugging:

```go
func deployWithLogging(ctx context.Context, client *sevalla.Client, appID string) error {
    log.Printf("Starting deployment for application: %s", appID)
    
    start := time.Now()
    deployment, resp, err := client.Applications.Deploy(ctx, appID)
    if err != nil {
        log.Printf("Deployment failed after %v: %v", time.Since(start), err)
        return err
    }
    
    log.Printf("Deployment started: %s (Status: %d)", deployment.ID, resp.StatusCode)
    return nil
}
```

## Error Handling

The SDK provides comprehensive error handling with helper functions:
}
```

```go
createReq := &sevalla.CreateApplicationRequest{
    Name:          "my-app",
    RepositoryURL: "https://github.com/user/repo",
    Branch:        "main",
    Region:      sevalla.RegionUSCentral,
    Plan:       sevalla.PlanStarter,
    EnvironmentVars: map[string]string{
        "NODE_ENV": "production",
        "API_KEY":  "secret",
    },
}

app, _, err := client.Applications.Create(ctx, createReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created app: %s with ID: %s\n", app.Name, app.ID)
```

#### Deploy an Application

```go
// Trigger a deployment
deployment, _, err := client.Applications.Deploy(ctx, app.ID, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deployment started: %s\n", deployment.ID)

// Wait for deployment to complete
finalDeployment, err := client.Deployments.WaitForDeployment(ctx, deployment.ID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deployment completed: %s\n", finalDeployment.State)
```

#### Scale an Application

```go
// Scale to 3 replicas
scaledApp, _, err := client.Applications.Scale(ctx, app.ID, 3)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Scaled to %d replicas\n", scaledApp.Replicas)
```

#### Update Environment Variables

```go
vars := map[string]string{
    "NEW_VAR": "value",
    "API_KEY": "updated-secret",
}

_, err = client.Applications.SetEnvironmentVariables(ctx, app.ID, vars)
if err != nil {
    log.Fatal(err)
}
```

### Databases

#### Create a Database

```go
createDBReq := &sevalla.CreateDatabaseRequest{
    Name:     "my-database",
    Type:     sevalla.EnginePostgreSQL,
    Version:  "14",
    Region: sevalla.RegionUSCentral,
    Size:     "small",
    Storage:  10, // GB
    Backups:  true,
}

db, _, err := client.Databases.Create(ctx, createDBReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Database created: %s\n", db.InternalURL)
```

#### List Databases

```go
databases, _, err := client.Databases.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, db := range databases {
    fmt.Printf("Database: %s (Type: %s, Size: %s)\n", 
        db.Name, db.Type, db.Size)
}
```

### Static Sites

#### Create a Static Site

```go
createSiteReq := &sevalla.CreateStaticSiteRequest{
    Name:            "my-site",
    RepositoryURL:   "https://github.com/user/site",
    Branch:          "main",
    BuildCommand:    "npm run build",
    OutputDirectory: "dist",
    Region:        sevalla.RegionEuropeWest,
}

site, _, err := client.StaticSites.Create(ctx, createSiteReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Static site created: %s\n", site.URL)
```

### Pipelines

#### Create a CI/CD Pipeline

```go
pipelineReq := &sevalla.CreatePipelineRequest{
    Name:    "test-and-deploy",
    Enabled: true,
    Trigger: "push",
    Branch:  "main",
    Steps: []sevalla.PipelineStep{
        {
            Name:    "test",
            Command: "npm test",
            Timeout: 300,
        },
        {
            Name:      "deploy",
            Command:   "npm run deploy",
            DependsOn: []string{"test"},
        },
    },
}

pipeline, _, err := client.Pipelines.Create(ctx, pipelineReq)
if err != nil {
    log.Fatal(err)
}

// Trigger a pipeline run
run, _, err := client.Pipelines.Run(ctx, pipeline.ID)
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

The SDK provides comprehensive error handling with helper functions:

```go
app, resp, err := client.Applications.Get(ctx, "app-123")
if err != nil {
    // Use helper functions to check error types
    switch {
    case sevalla.IsNotFound(err):
        fmt.Println("Application not found")
    case sevalla.IsUnauthorized(err):
        fmt.Println("Invalid API key or insufficient permissions")
    case sevalla.IsForbidden(err):
        fmt.Println("Access forbidden")
    case sevalla.IsConflict(err):
        fmt.Println("Resource conflict (e.g., duplicate name)")
    case sevalla.IsRateLimited(err):
        fmt.Println("Rate limit exceeded - retry later")
    case sevalla.IsServerError(err):
        fmt.Println("Server error - retry recommended")
    case sevalla.IsClientError(err):
        fmt.Println("Client error - check your request")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

### Available Error Helpers

- `IsNotFound(err)` - Resource not found (404)
- `IsUnauthorized(err)` - Authentication failed (401)
- `IsForbidden(err)` - Permission denied (403)
- `IsConflict(err)` - Resource conflict (409)
- `IsRateLimited(err)` - Rate limit exceeded (429)
- `IsClientError(err)` - Any 4xx error
- `IsServerError(err)` - Any 5xx error

### Getting Error Details

```go
if err != nil {
    // Extract APIError for detailed information
    if apiErr, ok := err.(*sevalla.APIError); ok {
        fmt.Printf("Status: %d\n", apiErr.StatusCode)
        fmt.Printf("Message: %s\n", apiErr.Message)
        fmt.Printf("Details: %v\n", apiErr.Details)
    }
}
```

## Response Handling

All API methods return three values: the resource, the response, and an error.

```go
app, resp, err := client.Applications.Get(ctx, "app-123")
if err != nil {
    log.Fatal(err)
}

// Access response metadata
fmt.Printf("Status Code: %d\n", resp.StatusCode)
fmt.Printf("Rate Limit: %d/%d\n", resp.RateLimit.Remaining, resp.RateLimit.Limit)
fmt.Printf("Rate Reset: %s\n", time.Until(resp.RateLimit.Reset))

// Pagination information
if resp.TotalCount > 0 {
    fmt.Printf("Total: %d\n", resp.TotalCount)
    fmt.Printf("Page: %d/%d\n", resp.CurrentPage, resp.TotalPages)
    fmt.Printf("Next Page: %d\n", resp.NextPage)
}
```

## Pagination

Use `ListOptions` for paginated results:

```go
opts := &sevalla.ListOptions{
    Page:    1,
    PerPage: 100,  // Maximum items per page
    Sort:    "created_at",
    Order:   "desc",
}

var allApps []*sevalla.Application

for {
    apps, resp, err := client.Applications.List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process apps
    allApps = append(allApps, apps...)
    fmt.Printf("Fetched %d apps (Page %d/%d)\n", len(apps), resp.CurrentPage, resp.TotalPages)
    
    // Check for more pages
    if resp.NextPage == 0 {
        break
    }
    opts.Page = resp.NextPage
}

fmt.Printf("Total applications: %d\n", len(allApps))
```

## Context Management

All API methods accept a `context.Context` for timeout and cancellation control:

```go
// Simple timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

app, _, err := client.Applications.Create(ctx, createReq)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Request timed out")
    } else {
        log.Fatal(err)
    }
}

// User cancellation
ctx, cancel := context.WithCancel(request.Context())
defer cancel()

// Long-running operation
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

deployment, _, err := client.Applications.Deploy(ctx, appID)
```

## Advanced Configuration

### Custom HTTP Client

Customize the underlying HTTP client for specific requirements:

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        // Add proxy if needed
        Proxy: http.ProxyFromEnvironment,
    },
}

client := sevalla.NewClient(
    sevalla.WithAPIKey("your-api-key"),
    sevalla.WithHTTPClient(httpClient),
)
```

### Custom Base URL

Use custom base URLs for testing or regional endpoints:

```go
// Staging environment
client := sevalla.NewClient(
    sevalla.WithAPIKey("your-staging-key"),
    sevalla.WithBaseURL("https://staging-api.sevalla.com/v2"),
)

// Local development
client := sevalla.NewClient(
    sevalla.WithAPIKey("dev-key"),
    sevalla.WithBaseURL("http://localhost:8080/v2"),
)
```

### Request/Response Debugging

Add middleware to log requests and responses:

```go
type debugTransport struct {
    wrapped http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // Log request
    log.Printf("→ %s %s", req.Method, req.URL)
    
    // Make request
    resp, err := t.wrapped.RoundTrip(req)
    if err != nil {
        log.Printf("✗ Request failed: %v", err)
        return nil, err
    }
    
    // Log response
    log.Printf("← %d %s", resp.StatusCode, resp.Status)
    return resp, nil
}

httpClient := &http.Client{
    Transport: &debugTransport{
        wrapped: http.DefaultTransport,
    },
}

client := sevalla.NewClient(
    sevalla.WithAPIKey("your-api-key"),
    sevalla.WithHTTPClient(httpClient),
)
```

## Complete Examples

### Example 1: Deploy Application with Monitoring

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
    // Initialize client
    client := sevalla.NewClient(
        sevalla.WithAPIKey(os.Getenv("SEVALLA_API_KEY")),
    )

    ctx := context.Background()

    // Create application
    createReq := &sevalla.CreateApplicationRequest{
        Name:          "production-api",
        RepositoryURL: "https://github.com/user/api",
        Branch:        "main",
        Region:      sevalla.RegionUSCentral,
        Plan:       sevalla.PlanPro,
        Replicas:      3,
        EnvironmentVars: map[string]string{
            "NODE_ENV":     "production",
            "DATABASE_URL": "postgresql://...",
        },
    }

    app, _, err := client.Applications.Create(ctx, createReq)
    if err != nil {
        log.Fatalf("Failed to create application: %v", err)
    }

    fmt.Printf("✓ Created application: %s (ID: %s)\n", app.Name, app.ID)

    // Trigger deployment
    deployment, _, err := client.Applications.Deploy(ctx, app.ID)
    if err != nil {
        log.Fatalf("Failed to start deployment: %v", err)
    }

    fmt.Printf("✓ Deployment started: %s\n", deployment.ID)

    // Monitor deployment
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    timeout := time.After(15 * time.Minute)

    for {
        select {
        case <-ticker.C:
            dep, _, err := client.Deployments.Get(ctx, deployment.ID)
            if err != nil {
                log.Fatalf("Failed to check deployment: %v", err)
            }

            fmt.Printf("⏳ Deployment state: %s\n", dep.State)

            switch dep.State {
            case sevalla.StatusSuccess:
                fmt.Println("✓ Deployment successful!")
                
                // Get application URL
                updatedApp, _, _ := client.Applications.Get(ctx, app.ID)
                fmt.Printf("✓ Application URL: %s\n", updatedApp.URL)
                return

            case sevalla.StatusFailed:
                // Get logs for debugging
                logs, _, _ := client.Deployments.GetLogs(ctx, deployment.ID)
                fmt.Printf("✗ Deployment failed!\n%s\n", logs)
                os.Exit(1)
            }

        case <-timeout:
            fmt.Println("✗ Deployment timed out")
            os.Exit(1)
        }
    }
}
```

### Example 2: Database Backup and Restore

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
    client := sevalla.NewClient(
        sevalla.WithAPIKey(os.Getenv("SEVALLA_API_KEY")),
    )

    ctx := context.Background()
    dbID := "db-123"

    // Create backup before maintenance
    fmt.Println("Creating backup...")
    backupReq := &sevalla.CreateBackupRequest{
        Type:        "manual",
        Description: "Pre-maintenance backup",
    }

    backup, _, err := client.Databases.CreateBackup(ctx, dbID, backupReq)
    if err != nil {
        log.Fatalf("Failed to create backup: %v", err)
    }

    fmt.Printf("✓ Backup created: %s\n", backup.ID)

    // Wait for backup to complete
    for i := 0; i < 30; i++ {
        b, _, err := client.Databases.GetBackup(ctx, dbID, backup.ID)
        if err != nil {
            log.Fatalf("Failed to check backup: %v", err)
        }

        if b.Status == "completed" {
            fmt.Printf("✓ Backup completed (Size: %d MB)\n", b.Size/1024/1024)
            break
        }

        fmt.Printf("⏳ Backup in progress: %s\n", b.Status)
        time.Sleep(10 * time.Second)
    }

    // Perform maintenance
    fmt.Println("Performing maintenance...")
    // ... maintenance operations ...

    // If something goes wrong, restore from backup
    shouldRestore := false  // Set based on your logic

    if shouldRestore {
        fmt.Println("Restoring from backup...")
        restoreReq := &sevalla.RestoreBackupRequest{
            BackupID: backup.ID,
        }

        _, err := client.Databases.RestoreFromBackup(ctx, dbID, restoreReq)
        if err != nil {
            log.Fatalf("Failed to restore: %v", err)
        }

        fmt.Println("✓ Restore initiated")
    }
}
```

### Example 3: Multi-Environment Pipeline

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    sevalla "github.com/juststeveking/sevalla-go"
)

func main() {
    client := sevalla.NewClient(
        sevalla.WithAPIKey(os.Getenv("SEVALLA_API_KEY")),
    )

    ctx := context.Background()

    // Create production pipeline
    pipelineReq := &sevalla.CreatePipelineRequest{
        Name:    "production-pipeline",
        Enabled: true,
        Trigger: "push",
        Branch:  "main",
        Steps: []sevalla.PipelineStep{
            {
                Name:    "install",
                Command: "npm ci",
                Image:   "node:18",
                Timeout: 300,
            },
            {
                Name:      "lint",
                Command:   "npm run lint",
                DependsOn: []string{"install"},
                Timeout:   180,
            },
            {
                Name:      "test",
                Command:   "npm test -- --coverage",
                DependsOn: []string{"install"},
                Timeout:   600,
            },
            {
                Name:      "build",
                Command:   "npm run build",
                DependsOn: []string{"lint", "test"},
                Timeout:   600,
            },
            {
                Name:      "deploy-staging",
                Command:   "npm run deploy:staging",
                DependsOn: []string{"build"},
                Timeout:   300,
            },
            {
                Name:      "e2e-tests",
                Command:   "npm run test:e2e",
                DependsOn: []string{"deploy-staging"},
                Timeout:   900,
            },
            {
                Name:      "deploy-production",
                Command:   "npm run deploy:production",
                DependsOn: []string{"e2e-tests"},
                Timeout:   300,
            },
        },
        Environment: map[string]string{
            "NODE_ENV": "production",
            "CI":       "true",
        },
    }

    pipeline, _, err := client.Pipelines.Create(ctx, pipelineReq)
    if err != nil {
        log.Fatalf("Failed to create pipeline: %v", err)
    }

    fmt.Printf("✓ Pipeline created: %s (ID: %s)\n", pipeline.Name, pipeline.ID)

    // Trigger pipeline run
    run, _, err := client.Pipelines.Run(ctx, pipeline.ID)
    if err != nil {
        log.Fatalf("Failed to trigger pipeline: %v", err)
    }

    fmt.Printf("✓ Pipeline run started: %s\n", run.ID)
    fmt.Printf("→ View logs at: https://sevalla.com/pipelines/%s/runs/%s\n", pipeline.ID, run.ID)
}
```

## Available Services

- **Applications** - Complete application lifecycle management (create, deploy, scale, logs, custom domains, etc.)
- **Databases** - Database management for PostgreSQL, MySQL, MongoDB, and Redis (create, credentials, backups, restore, usage)
- **StaticSites** - Static site deployment (frameworks: Next.js, Vite, Hugo, Jekyll, etc.)
- **Deployments** - Deployment monitoring and control (status, logs, cancel)
- **Pipelines** - CI/CD pipeline automation (create, run, monitor, retry)
- **StaticSites** - Manage static sites
- **Deployments** - Monitor and control deployments
- **Pipelines** - CI/CD pipeline management

## Available Types

### Regions
- `RegionUSCentral` - US Central region
- `RegionUSEast` - US East region  
- `RegionEuropeWest` - Europe West region
- `RegionAsiaSouth` - Asia South region

### Plans
- `PlanHobby` - Hobby tier
- `PlanStarter` - Starter tier
- `PlanPro` - Pro tier
- `PlanBusiness` - Business tier
- `PlanEnterprise` - Enterprise tier

### Database Engines
- `EnginePostgreSQL` - PostgreSQL database
- `EngineMySQL` - MySQL database
- `EngineMongoDB` - MongoDB database
- `EngineRedis` - Redis cache

## Testing

Create a test file to verify the SDK:

```go
// main_test.go
package main

import (
    "context"
    "os"
    "testing"
    
    sevalla "github.com/juststeveking/sevalla-go"
)

func TestSDK(t *testing.T) {
    apiKey := os.Getenv("SEVALLA_API_KEY")
    if apiKey == "" {
        t.Skip("SEVALLA_API_KEY not set")
    }
    
    client := sevalla.NewClient(
        sevalla.WithAPIKey(apiKey),
    )
    
    ctx := context.Background()
    
    // Test listing applications
    apps, resp, err := client.Applications.List(ctx, nil)
    if err != nil {
        t.Fatalf("Failed to list applications: %v", err)
    }
    
    t.Logf("Found %d applications", len(apps))
    t.Logf("Response status: %d", resp.StatusCode)
}
```

Run tests:
```bash
export SEVALLA_API_KEY="your-api-key"
go test -v
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Next Steps

Now that you have the SDK ready:

1. **Test the SDK** with your actual Sevalla API key
2. **Build the MCP Server** using this SDK (much cleaner!)
3. **Publish to GitHub** and tag as v0.1.0
4. **Add more features** as needed

The SDK is now ready to be used by your MCP server, making the implementation much cleaner and type-safe!
