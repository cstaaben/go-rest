# PRD-0005: Request History

## Problem Statement

When building or debugging APIs, I often need to compare the outcomes of different requests over time, review what parameters were sent in past runs, or restore a previous successful response. Currently, there is no log of executions, making it impossible to see past requests or track down when an API behavior changed.

## Solution

Implement a Request History feature. Every execution of a Saved Request will write a detailed execution snapshot to a rolling append-only JSON Lines file (`history.jsonl`). Integrate this history into the GUI via a new **History** tab in the sidebar, allowing the user to browse past executions, inspect what was sent, and load previous request/response states.

## User Stories

1. As a user, I want my request runs to be logged automatically when I execute them, so that I have a record of my interactions with the API.
2. As a user, I want only executions of Saved Requests (including ad-hoc requests executed with the `--save` flag) to be logged in history, so that throwaway, unsaved requests (which may contain transient credentials) are kept out of the persistent log.
3. As a user, I want each history log entry to record a complete **Request Snapshot** (the exact URL, headers, and body sent *after* environment variable substitution), along with the response details (status, headers, body), timestamp, and duration, so that I can see exactly what went over the wire.
4. As a user, I want the history log size to be automatically capped using a rolling limit (defaulting to the last 1000 runs) that I can customize via the `history_limit` setting in my configuration file, so that the log doesn't consume excessive disk space or slow down the application.
5. As a user, I want to access my history via a **History** tab in the sidebar alongside my **Requests** collection, so that I can easily scan past executions.
6. As a user, I want to click a history entry in the sidebar, so that its historical snapshot payload and response details load into my request and response editors for review or re-run.

## Implementation Decisions

* **Log Storage & Format:**
  * File name: `history.jsonl` under the `data_dir` directory.
  * Format: JSON Lines (each line is a single JSON object).
  * Fields per entry:
    * `request_id`: UUID of the saved request
    * `timestamp`: ISO-8601 execution time
    * `duration_ms`: Execution time in milliseconds
    * `request_snapshot`: `{ url, method, headers, body }` containing substituted values
    * `response_snapshot`: `{ status_code, headers, body }`
* **Rolling Pruning:**
  * Every time a history entry is successfully appended, the history manager checks the line count.
  * If the count exceeds the configured `history_limit` (default: 1000), it prunes the oldest lines from the beginning of the file.
* **GUI Layout:**
  * The GUI sidebar is updated from a simple Tree container to a tabbed container with two tabs: **Requests** and **History**.
  * The History tab displays a list of past runs sorted chronologically (latest first).

## Testing Decisions

* **Testing Philosophy:**
  * Verify the logging and rolling pruning behavior using unit tests. Mock the filesystem path to a temporary folder and mock the `history_limit` to a small number (e.g. 5) to verify that appending works, transient requests are ignored, and older entries are pruned correctly.
* **Modules to Test:**
  * `internal/history` (the history manager and file writer).
* **Prior Art:**
  * None.

## Out of Scope

* Syncing history files across multiple environments or machines.
* Exporting history entries as independent files (e.g. exporting a history item to a curl script — this is covered by the snippet generation feature).
