# PRD-0006: Code Snippet Generation

## Problem Statement

When I need to migrate a request that works in `go-rest` into my application code or share it with a colleague who doesn't use the client, I have to manually translate the URL, headers, and request body into code syntax. This is tedious, error-prone, and slow.

## Solution

Implement a Code Snippet Generation engine. The engine will compile a Request definition into equivalent code for four core targets: `curl` (shell), Go (`net/http`), Python (`requests`), and JavaScript (`fetch`). Provide a header action button in the GUI to view and copy the snippet, and a `go-rest snippet` subcommand in the CLI to output the snippet directly to the terminal.

## User Stories

1. As a developer, I want to generate a code snippet from any Request, so that I can quickly integrate it into my external scripts.
2. As a developer, I want to choose between `curl`, Go (`net/http`), Python (`requests`), and JavaScript (`fetch`) target formats, so that I get the correct syntax for my development stack.
3. As a developer, I want to toggle between generating the code snippet using raw placeholders (preserving `{{variable}}` templates) or fully resolved active environment values, so that I can copy either a copy-pasteable script or a generic code template.
4. As a developer, I want to access this feature in the GUI by clicking a dedicated code button (`</>`) in the main header (next to the "Send" button), which opens a modal with the generated code, a "Copy to Clipboard" button, and configuration controls.
5. As a developer, I want to generate code directly from the command line using `go-rest snippet <request-identifier>`, so that I can view and pipe the generated code snippet to other terminal tools.
6. As a developer, I want the CLI `snippet` command to accept `--lang`/`-l` (defaulting to `curl`) and `--raw` (to skip variable resolution), so that I can control code output easily.

## Implementation Decisions

* **Snippet Compiler:**
  * Implement a stateless compiler module that takes a Request (or Request Snapshot) and formats it into the target language's code structure.
  * Targets:
    * `curl`: Constructs a single `curl` CLI command with correct flags (`-X`, `-H`, `-d`).
    * Go: Generates code using the standard library's `net/http` package (creating client, request, headers, body reader).
    * Python: Generates code using the `requests` library.
    * JavaScript: Generates code using the native `fetch` API.
* **Variable Resolution Toggle:**
  * The GUI dialog and CLI parser pass either the raw Request definition (with placeholders intact) or the resolved Request Snapshot to the snippet compiler depending on the user's choice.
* **CLI Command:**
  * Command: `go-rest snippet <request-id-or-path>`.
  * Outputs the code snippet directly to `stdout`.

## Testing Decisions

* **Testing Philosophy:**
  * Test the compiler package with stateless unit tests. Pass a variety of requests (having empty body, JSON body, multiline headers, different HTTP methods, and variables) and assert that the generated string matches the expected syntax templates (both raw and resolved).
* **Modules to Test:**
  * `internal/snippet` (the code generator module).
* **Prior Art:**
  * None.

## Out of Scope

* Supporting additional frameworks (e.g. Axios, HTTPie, PHP cURL) beyond the four core targets.
* Syntactic pretty-printing or linting of the generated Go/Python/JS code (output standard, valid code format).
