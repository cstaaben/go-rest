# HTTP Request Execution

We decided to structure request execution as follows:
1. The core client wrapper (`internal/client`) exposes a synchronous execution method:
   ```go
   func (c *Client) Send(ctx context.Context, req *request.Request) (*response.Data, error)
   ```
2. The response data is returned in a dedicated `response.Data` struct containing the HTTP headers, response body, protocol version, status code, and execution duration.
3. The execution method accepts a Go `context.Context` to allow cancellation of requests in flight and timeout control by the callers (GUI or CLI).
