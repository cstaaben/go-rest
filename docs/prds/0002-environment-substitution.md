# PRD-0002: Environment Variable Substitution

## Problem Statement

When testing APIs across different stages (such as development, staging, and production), I am forced to duplicate my request definitions or manually edit parameters like host names, headers, and authorization tokens. This is tedious, error-prone, and clutters the request collection.

## Solution

Implement an Environment variable substitution engine. Allow users to reference variables in URLs, header keys, header values, and request bodies using the `{{variable_name}}` syntax. The engine will dynamically replace placeholders with values from the selected Active Environment at the moment of execution.

## User Stories

1. As a developer, I want to use `{{variable_name}}` placeholders in a Request's URL, headers, and body, so that I don't have to hardcode configuration parameters.
2. As a developer, I want to select an Active Environment, so that all requests automatically resolve variable placeholders using that environment's variables.
3. As a developer, I want request execution to fail immediately with a clear validation error if a placeholder cannot be resolved, so that I do not send malformed HTTP requests (e.g., calling `DELETE /api/items/` instead of `DELETE /api/items/123`).
4. As a developer, I want the environment substitution to have no fallback to operating system environment variables, so that request execution is strictly reproducible across different developer machines.
5. As a developer, I want the configuration parser to validate that the selected active environment name matches an existing environment file on startup.

## Implementation Decisions

* **Placeholder Syntax:**
  * Double curly braces wrapper: `{{variable_name}}`.
* **Resolution Scope:**
  * Substitution applies to the Request URL, HTTP method (optional but URL and headers are mandatory), all Header keys and values, and the Request Body.
  * Extracted Session Variables (if any) take precedence over the Active Environment's variables when resolving placeholders.
* **Validation & Failure Safety:**
  * If a request contains any `{{placeholder}}` that cannot be resolved in the Active Environment (or Session Variables), the engine halts execution immediately and returns a detailed validation error listing the missing variables.
  * System environment variables (`os.Getenv`) are explicitly ignored.

## Testing Decisions

* **Testing Philosophy:**
  * The substitution engine should be tested as a pure, stateless string processor. Tests should assert exact replacement strings and check that execution halts with the expected error list when variables are missing.
* **Modules to Test:**
  * `internal/environment` (specifically the substitution resolver function).
* **Prior Art:**
  * None currently in codebase. We will introduce a new test file, e.g., `internal/environment/substitution_test.go`.

## Out of Scope

* Dynamic evaluation functions inside placeholders (e.g., `{{$randomInt}}` or `{{$timestamp}}`).
* Nested variable resolution (e.g., resolving `{{user_{{id}}}}`).

## Further Notes

* Supported variable types are primitive scalars (strings, numbers, booleans) which are cast to strings during substitution.
