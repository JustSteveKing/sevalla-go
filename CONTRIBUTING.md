# Contributing

Thanks for your interest in contributing! This repo is a Go SDK for the Sevalla API.

## Prerequisites

- Go 1.23+
- Make

## Setup

1. Fork and clone the repo
2. Install tools and deps:
   - `make install`
3. Run the full suite:
   - `make pre-commit`

## Common Tasks

- Format: `make fmt`
- Lint: `make lint`
- Test: `make test`
- Security scan: `make security`
- Generate docs: `make docs`

## Branching & PRs

- Create feature branches from `main`
- Ensure CI is green and coverage not reduced
- Update README/docs where applicable

## Releases

Maintainers can tag a release via `make release` which creates and pushes a tag.
