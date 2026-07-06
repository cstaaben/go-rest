# Request History

We decided to implement a Request History feature to track executions of Saved Requests:
1. History is stored in a single append-only JSON Lines file (`history.jsonl`) under the `data_dir`.
2. Only executions of **Saved Requests** (including ad-hoc requests run with the `--save` flag) are logged. Transient or unsaved requests are excluded from history.
3. Each entry in the log records a **Request Snapshot** containing the exact URL, headers, and body sent (after variable substitution), along with the response status, body, headers, duration, and execution timestamp.
4. The history file has an automatic rolling cap to prevent unbounded growth. The limit defaults to the last 1000 executions and is user-configurable via the `history_limit` setting in the main configuration file (`go-rest.yaml`).
5. The GUI displays history in a **History** sidebar tab alongside the **Requests** collection tab. Selecting an entry loads the snapshot and response into the active workspace.
