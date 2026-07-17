# AGENTS Instructions

## 1. Project Overview

`go-rest` is a Terminal User Interface (TUI) command-line REST client.

- **TUI Framework**: Built using [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), and [Lipgloss](https://github.com/charmbracelet/lipgloss).
- **Architecture**:
  - `cmd/go-rest`: Entry point and command line flag/config parsing.
  - `internal/ui`: Contains the UI components (request list, request editor, response viewer, environment editor).
  - `internal/client`: Handles request execution.
  - `internal/model`: Global state and data structures (targets, configurations).

## 2. Build and Testing Commands

We use **Mage** for workflow automation, but standard `go` CLI tools work as well.

- **Build**:
  - Mage: `mage build`
  - Go: `go build -o bin/ ./cmd/go-rest`
- **Run**:
  - Mage: `mage run` (Runs example configuration)
  - Go: `go run ./cmd/go-rest --config ./example/config.yaml`
- **Lint**:
  - Mage: `mage lint`
  - Go/Shell: `golangci-lint run --allow-parallel-runners`
- **Tidy Modules**:
  - Mage: `mage tidy`
- **Vendor Modules**:
  - Mage: `mage vendor`
- **Kill Binary**:
  - Mage: `mage kill`
- **Test**:
  - Go: `go test ./...`

## 3. Code Style

- **Conventions**: Write compact, readable Go code. Use latest stable features.
- **Strict Typing**: Avoid empty interfaces (`interface{}`) and `any` types where possible. Always define clear Go structs.
- **Variable Names**: Use verbose, descriptive variable names (e.g., `isUserAuthenticated` instead of `auth`).
- **Documentation**: Provide Go docstrings/comments documenting all new packages, structs, and public functions. Update comments immediately if you modify the underlying logic.
- **UI Design**: Maintain bubbletea MVC architecture. Keep subcomponent layout updates reactive based on window size messages (`tea.WindowSizeMsg`).

## 4. Testing Instructions

- **Writing Tests**: When fixing a bug, first write a reproducing unit test in the relevant package, then fix the bug.
- **Running Tests**: Run `go test ./...`. Ensure all tests compile and pass.
- **Mocking**: Mock any external network requests or complex filesystem reads to keep unit tests fast and deterministic.

## 5. Security Considerations

- **Secrets**: NEVER print, log, or commit API keys, authorization tokens, or sensitive credentials.
- **Config & Logs**: Do not print authorization headers or query parameters containing sensitive tokens into standard logs. Output warnings/info logs to `./logs` directory without sensitive payloads.
- **Network Safety**: Ensure REST requests are sent only to the hostnames explicitly configured by the user. Do not make unauthorized background HTTP calls.

## 6. Agent Skills & Workflows

### Issue tracker

Issues are tracked in GitHub Issues. External PRs are not treated as a triage surface. See [issue-tracker.md](file:///home/corbinstaaben/code/src/github.com/cstaaben/go-rest-agent-setup/docs/agents/issue-tracker.md).

### Triage labels

Default triage labels are used (each role maps directly to the label string). See [triage-labels.md](file:///home/corbinstaaben/code/src/github.com/cstaaben/go-rest-agent-setup/docs/agents/triage-labels.md).

### Domain docs

Single-context repo layout with a single `CONTEXT.md` and `docs/adr/` directory at the repo root. See [domain.md](file:///home/corbinstaaben/code/src/github.com/cstaaben/go-rest-agent-setup/docs/agents/domain.md).
