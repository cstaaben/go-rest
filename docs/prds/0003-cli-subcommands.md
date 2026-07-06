# PRD-0003: CLI Subcommands and Ad-hoc Execution

## Problem Statement

When I need to quickly run an API request in the terminal or automate it in a script (like checking an endpoint status or triggering a webhook), I am forced to launch the graphical user interface. This breaks terminal flow and makes `go-rest` impossible to integrate into scripts, CI pipelines, or cron tasks.

## Solution

Extend the `go-rest` CLI with a minimalist set of subcommands using the Cobra library. Support executing saved requests, listing collections, and building ad-hoc (one-off) requests directly from command-line flags. Ensure that execution outputs the response body by default to allow seamless piping into terminal utilities like `jq`.

## User Stories

1. As a developer, I want running the bare `go-rest` command with no arguments to launch the GUI, so that the main interactive mode remains frictionless.
2. As a developer, I want to execute a saved request from the CLI using `go-rest run <id-or-path>`, so that I can run requests from terminal history or scripts.
3. As a developer, I want to list all saved Groups and Requests in a clean text layout using `go-rest list`, so that I can find request identifiers without opening the GUI.
4. As a developer, I want to run one-off ad-hoc requests using flags (`--url`, `--method`, `--header`, `--body`), so that I don't have to define a YAML file for a simple test.
5. As a developer, I want to save an ad-hoc request to disk using the `--save` flag (with optional `--name` and `--group`), so that I can capture a successful experiment into my collection.
6. As a developer, I want variable placeholders in my ad-hoc requests to resolve against a specified environment (e.g. `go-rest run --url "{{baseUrl}}/users" --env dev`), failing immediately if no environment is selected or variables are missing.
7. As a developer, I want `go-rest run` to output only the raw response body by default, so that I can easily pipe it to other tools (e.g. `go-rest run my_req | jq .name`).
8. As a developer, I want formatting options (`--include`/`-i` for response headers, `--verbose`/`-v` for full headers, and `--quiet`/`-q` to suppress all output), so that I can control terminal noise during scripting.

## Implementation Decisions

* **Command Router:**
  * Integrate `github.com/spf13/cobra` to handle command routing.
  * If no subcommands or flags are passed (and `stdout` is a TTY), execute the default action (launch the GUI).
* **Command Syntax:**
  * `go-rest list`: Walks the requests directory and prints a flat or nested tree.
  * `go-rest run <id-or-path>`: Resolves, executes, and outputs the request.
  * Ad-hoc execution: `go-rest run --url <url> [--method <method>] [--header <header>] [--body <body>]`.
* **Output Handling:**
  * Default: Write only the raw response body to `stdout`.
  * `--include` / `-i`: Prepend response headers before the body.
  * `--verbose` / `-v`: Print outgoing request headers/body and incoming response headers/body.
  * `--quiet` / `-q`: Direct all response output to `/dev/null` (exit code reflects HTTP success/failure).

## Testing Decisions

* **Testing Philosophy:**
  * Test the command-line interface by invoking the command routing functions programmatically. Redirect standard output and error buffers to capture and verify string output, avoiding starting a true terminal environment.
* **Modules to Test:**
  * `cmd/go-rest` (Cobra command parser, argument validations, and output formatter).
* **Prior Art:**
  * None.

## Out of Scope

* CLI interactive prompts for editing or creating requests (users should edit YAML files or use the GUI).
* Interactive progress bars (keep CLI output clean for scripting).

## Further Notes

* Non-zero exit codes will be returned if:
  * The HTTP request fails (network error).
  * A validation error occurs (e.g. missing environment variables).
  * The server returns a client/server error status (HTTP 4xx or 5xx), unless a specific flag overrides this behavior.
