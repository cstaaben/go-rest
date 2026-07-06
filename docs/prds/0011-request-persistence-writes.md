# PRD-0011: Request Persistence Writes

## Problem Statement

When I modify request parameters, create new groups, or reorganize collections in the GUI, my changes are not written back to disk. Every time I exit the application, all additions, edits, deletions, and structural reorganizations are completely lost.

## Solution

Implement database write persistence functions in the `internal/request` package. Saving or deleting a Request will overwrite the parent Group's YAML file. Group creations, renames, or moves will update their specific metadata file. Deleting a Group will remove its file and recursively delete all child Group files associated with it on disk.

## User Stories

1. As a user, I want my changes to a Request (URL, method, headers, body, auth, variables) to be automatically saved, so that my modifications are persisted across application restarts.
2. As a user, I want to create a new Request within a Group, so that it is appended to that Group's YAML file on disk.
3. As a user, I want to create, rename, and move Groups in the UI and have these operations written back to disk, so that my collection structure is saved.
4. As a user, I want to delete a Group and have its YAML file removed, along with the recursive deletion of all its sub-groups and requests on disk, so that obsolete collection directories are completely cleaned up.
5. As a developer, I want all file write operations to be synchronous, so that my backend code can ensure data consistency immediately and run simple, race-free tests.

## Implementation Decisions

* **Save and Delete Functions:**
  * Add write operations to the `internal/request` package:
    * `SaveGroup(dataDir string, g *Group) error`: Marshals the group struct to YAML and overwrites `<group_id>.yaml` in the `/requests` subdirectory.
    * `DeleteGroup(dataDir string, groupID string) error`: Deletes `<group_id>.yaml` from disk.
* **Request Level Wrapper:**
  * Since Requests are stored inline, to edit or save a Request, the engine locates the parent Group file, updates the Request in the Group's in-memory array, and calls `SaveGroup`.
* **Recursive Deletion:**
  * When `DeleteGroup` is called, the engine scans the requests folder for any Group file whose `parent_id` matches the deleted Group's ID, and recursively calls `DeleteGroup` on them.

## Testing Decisions

* **Testing Philosophy:**
  * Test writes by writing to a temporary directory. Assert that file write actions successfully serialize to valid YAML, edits update the correct sub-keys in the file, and deleting a parent group cleans up all child files recursively.
* **Modules to Test:**
  * `internal/request` (YAML writers and directory crawlers).
* **Prior Art:**
  * Existing unit tests in `internal/request/load_test.go`.

## Out of Scope

* Automated git staging or committing upon writes (all operations strictly write to the local filesystem).
* File lock management or multi-process write conflict resolution (assume single-user concurrency).
