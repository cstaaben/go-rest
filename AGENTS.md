# Coding Agent Instructions

This document provides instructions for AI coding agents working on this project.

## Project Overview

- **Name:** go-rest
- **Type:** CLI and GUI
- **Stage:** Prototype

## Tech Stack

- **Language:** Go
- **Styling:** gofumpt
- **Linting:** golangci-lint
- **UI Components:** BubbleTea (TUI), Fyne (GUI)
- **Database:** None currently
- **Testing:** testify/assert
- **Package Manager:** Go modules

## Directory Structure

```text
├── assets            # GUI assets
├── cmd               # program entrypoint
├── example           # example configs/data
├── examples          # example programs for bubbletea and fyne frameworks
├── internal          
│   ├── client        # HTTP client wrapper
│   ├── config        # config types and parsing
│   ├── environment   # request environments
│   ├── model         # TUI model state
│   ├── request       # HTTP request wrapper
│   ├── ui            # TUI model definitions
│   └── uuid          # UUID utility
└── magefiles         # build tool
```

## Development Commands

```bash
# download dependencies
mage tidy vendor

# build
mage build

# lint
mage lint

# run tests
go test -shuffle=on ./...

# format code
gofumpt -w ./...
```

## Code Quality

- Maximum file length: 500 lines — split larger files into logically grouped packages
- No `fmt.Println` in production paths — use a structured logger
- All exported functions must have godoc comments
- Use `gofumpt` to format all code
- Use `golangci-lint` to lint all code
- Use context7 to find the most up-to-date Go documentation

## Testing Requirements

- Write unit tests for every new utility function
- Run `go test` before marking a task complete
- Put tests in a `*_test` package whenever possible

## Safety Guardrails (Critical)

- Never delete files without explicit user confirmation
- Never commit `.env`, `.env.local`, or any file containing secrets
- Never hardcode credentials, API keys, or tokens
- Validate all user input before processing

## Architecture Rules

- Do not export a function unless it is being used outside its package
- Whenever a function requires more than 4 arguments, use functional options or
  a settings struct
- **NEVER** ignore errors with `_`

## Communication

- Be concise — skip explanations of basic concepts
- When suggesting a change, explain the 'why', not just the 'what'
- If you notice a potential bug while working on something else, stop and flag it
- Always suggest the simplest solution that meets the requirements
- When uncertain about intent, ask rather than guessing
- Propose a plan for all changes before implementing them

## Agent skills

### Issue tracker

Issues are tracked using GitHub Issues (uses the `gh` CLI). Pull Requests are not treated as a request surface for triage. See `docs/agents/issue-tracker.md`.

### Triage labels

Using standard default labels: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repository layout (one `CONTEXT.md` + `docs/adr/` at the repo root). See `docs/agents/domain.md`.
