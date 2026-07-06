# PRD-0010: HTTP Request Execution

## Problem Statement

The REST client is currently just a visual shell; it has no underlying capability to send network requests, parse HTTP responses, or measure execution performance. I cannot interact with real APIs or see response payloads.

## Solution

Implement a synchronous, stateless HTTP request runner in the `internal/client` package. The runner will accept a `context.Context` (enabling timeout and cancellation control) and return a structured `response.Data` object enclosing the HTTP status code, response headers, response body, protocol version, and the execution duration in milliseconds.

## User Stories

1. As a user, I want to execute a configured Request, so that I can see the response returned by the target API server.
2. As a user, I want to view the time it took to receive a response (duration in milliseconds), so that I can monitor API performance.
3. As a user, I want to cancel an in-flight request immediately (e.g. by clicking a "Cancel" button in the GUI or sending Ctrl+C in the CLI), so that I don't have to wait for hung or long-running queries to time out.
4. As a developer, I want the core client runner to be synchronous and context-driven, so that calling frameworks (CLI and GUI) can manage their own threading/goroutine pools to prevent interface freezing.

## Implementation Decisions

* **Client Wrapper:**
  * Implement standard HTTP request dispatching in the `internal/client` package.
  * Method signature:
    ```go
    func (c *Client) Send(ctx context.Context, req *request.Request) (*response.Data, error)
    ```
* **Response Struct:**
  * Define `response.Data` inside a new package `internal/response`:
    ```go
    type Data struct {
        URL        string              `json:"url"`
        Method     string              `json:"method"`
        StatusCode int                 `json:"status_code"`
        Proto      string              `json:"proto"`
        Headers    map[string][]string `json:"headers"`
        Body       string              `json:"body"`
        DurationMs int64               `json:"duration_ms"`
    }
    ```
* **Timeout and Cancellation:**
  * Map incoming context to the `http.Request` using `req.WithContext(ctx)` before calling `http.Client.Do()`.
  * Track duration using standard `time.Now()` and `time.Since()`.

## Testing Decisions

* **Testing Philosophy:**
  * Test execution and cancellation by using Go's standard `net/http/httptest` package. Start a local mock HTTP server that simulates delays, custom headers, and body structures, asserting that:
    * The returned `response.Data` matches exactly.
    * Duration is recorded correctly.
    * Cancelling the context mid-request immediately aborts the call and returns the context cancellation error.
* **Modules to Test:**
  * `internal/client` (the request runner).
* **Prior Art:**
  * None.

## Out of Scope

* Importing custom SSL root certificate authorities (CAs) or configuring client-side keystores.
* Proxy configurations (HTTP/SOCKS proxy bypass).

## Further Notes

* Connection pooling and reuse will leverage the default configurations of Go's `http.DefaultTransport`.
