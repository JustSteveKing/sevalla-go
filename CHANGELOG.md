# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-10-18

### Added

- **CompanyID Filtering**: Add `CompanyID` field to `ListOptions` struct for filtering resources by company ID ([#701307e](https://github.com/JustSteveKing/sevalla-go/commit/701307e))
  - Enables filtering across all list operations (Applications, Databases, Static Sites, Deployments, Pipelines)
  - Optional field maintains full backward compatibility
  - Implemented as URL query parameter for consistent behavior
  - Usage example: `ListOptions{CompanyID: "company-123"}`

### Tests

- Add `TestStaticSitesService_ListWithCompanyID` test to verify company ID filtering ([#7b110a5](https://github.com/JustSteveKing/sevalla-go/commit/7b110a5))
  - Verifies query parameter is correctly passed to API
  - Maintains 79.7% code coverage

### Documentation

- Add complete working example in `examples/list-static-sites-by-company/` ([#54d4009](https://github.com/JustSteveKing/sevalla-go/commit/54d4009))
  - Demonstrates filtering static sites by company ID
  - Includes both environment variable and command-line argument support
- Update README.md and API.md with CompanyID usage examples ([#b61491f](https://github.com/JustSteveKing/sevalla-go/commit/b61491f))

### Breaking Changes

None - this release is fully backward compatible with v0.1.0.

## [0.1.0] - 2025-10-17

### Added

- Initial release of the Sevalla Go SDK
- Complete API coverage for all Sevalla endpoints:
  - Applications service (create, deploy, scale, manage)
  - Databases service (PostgreSQL, MySQL, MongoDB, Redis)
  - Static Sites service (deploy and manage static websites)
  - Deployments service (monitor and control deployments)
  - Pipelines service (CI/CD pipeline automation)
- Type-safe API with strongly typed requests and responses
- Comprehensive error handling with helper functions
- Built-in pagination support
- Full context.Context support for timeouts and cancellation
- 79.7% test coverage with comprehensive test suite
- Production-ready with minimal dependencies
- Complete documentation and examples

### Features

- Service-based client architecture
- Automatic request/response handling
- Rich error types (`IsNotFound`, `IsUnauthorized`, `IsRateLimited`, etc.)
- Helper functions for pointer types (`String`, `Int`, `Bool`)
- Pagination metadata in responses
- Custom HTTP client support
- Custom base URL support (for testing or regional endpoints)

[0.2.0]: https://github.com/JustSteveKing/sevalla-go/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/JustSteveKing/sevalla-go/releases/tag/v0.1.0
