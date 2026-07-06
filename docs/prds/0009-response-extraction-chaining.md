# PRD-0009: Response Extraction and Request Chaining

## Problem Statement

When testing multi-step API flows (such as running a POST request to login, extracting an access token from the response, and then using that token to authenticate subsequent requests), I must manually copy the token from the response panel and paste it into my environment file or request headers. This is slow, repetitive, and interrupts workflow.

## Solution

Implement a declarative Response Extraction and Request Chaining engine in the `internal/variables` package. Allow users to define **Extraction Rules** to capture values from HTTP response bodies (JSON using JSONPath, XML using XPath) and response headers, binding them to **Session Variables**. The session variables are stored in a temporary, in-memory store that overrides environment variables for the current session. Configure these rules via a new **Variables** tab in the GUI Request Editor, and via `--extract` flags in the CLI.

## User Stories

1. As a developer, I want to define post-request extraction rules on a Request, so that values are automatically extracted from the response upon successful execution.
2. As a developer, I want to extract values from JSON response bodies using **JSONPath** expressions (e.g., `$.access_token`), so that I can easily parse JSON APIs.
3. As a developer, I want to extract values from XML response bodies using **XPath** expressions (e.g., `/Envelope/Body/Token`), so that I can support XML and legacy SOAP APIs.
4. As a developer, I want the XML parser to automatically compile simple slash-delimited paths into namespace-agnostic XPath queries (using `*[local-name()='...']`) behind the scenes, so that I don't have to manually configure XML namespace prefixes.
5. As a developer, I want to extract response headers by name (e.g. `Set-Cookie` or `Authorization`), so that I can capture headers.
6. As a developer, I want extracted values to be saved as **Session Variables** in a temporary, in-memory store that overrides Active Environment variables during the active session, so that short-lived tokens do not dirty my permanent environment files on disk (avoiding Git configuration noise).
7. As a developer, I want to define extraction rules in the GUI Request Editor under a new **Variables** tab (alongside Body, Headers, Auth, Cookies), and specify them in the CLI via `--extract` flags (e.g., `--extract "token=body:$.access_token"`), so that chaining is easy to set up.

## Implementation Decisions

* **Extraction Rule Schema:**
  * Each extraction rule consists of:
    * `variable_name`: Target session variable name (e.g., `sessionToken`).
    * `source`: `body` | `header`.
    * `query`: The JSONPath expression, XPath expression, or header key name.
    * `type`: `jsonpath` | `xpath` | `header` (inferred from format/source or declared explicitly).
* **XPath Compiler:**
  * If the type is `xpath` and the query is a simple path (contains `/` and alphanumeric segments), the parser compiles it to a namespace-agnostic XPath query. E.g., `/Envelope/Body/Token` becomes `/*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='Token']`.
* **Session Store:**
  * Define a thread-safe, in-memory key-value map for Session Variables during the application session.
  * When executing requests, the variable resolver checks the Session Store map first, falling back to the Active Environment.
  * The session map is cleared when the application exits.

## Testing Decisions

* **Testing Philosophy:**
  * Verify value extraction and query formatting using unit tests. Mock JSON/XML response bodies and headers, pass different extraction rules, and assert that:
    * Extracted values match expected outputs.
    * The XPath compiler correctly generates and evaluates namespace-agnostic queries on namespaced XML payloads.
    * In-memory Session Variables successfully override Active Environment variables during variable resolution.
* **Modules to Test:**
  * `internal/variables` (the extraction resolver and XPath compiler).
* **Prior Art:**
  * None.

## Out of Scope

* Permanent write-back of extracted variables to the configuration YAML files on disk.
* Chaining executions across different collections (variables are scoped globally to the active application session).
