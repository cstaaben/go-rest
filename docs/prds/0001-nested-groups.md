# PRD-0001: Nested Groups

## Problem Statement

As my REST request collection grows, I cannot organize my requests into a hierarchical folder structure. Having all request groups in a flat list makes the sidebar cluttered, hard to navigate, and difficult to manage for larger projects.

## Solution

Support nested Request Groups. Allow a Group to reference a parent Group to form a hierarchical tree. Reconstruct this tree structure in memory for both GUI tree rendering and CLI listing, while keeping the physical filesystem storage flat to keep files small and version control diffs clean.

## User Stories

1. As a developer, I want to create a sub-group within an existing Group, so that I can organize related sub-resources under a parent resource (e.g., `/users/auth` under `/users`).
2. As a developer, I want to nest groups multiple levels deep, so that I can represent complex nested API hierarchies.
3. As a developer, I want to move a Group to be a child of another Group or move it back to the root level, so that I can reorganize my collections.
4. As a developer, I want to view my requests in a hierarchical tree structure in the GUI sidebar, so that I can easily scan and navigate my API endpoints.
5. As a developer, I want to delete a Group, which also deletes all of its nested child Groups and Requests, so that I don't leave orphaned data behind.
6. As a developer, I want to store my groups as flat YAML files in the filesystem with `parent_id` pointers, so that my local workspace is simple to inspect and git diffs remain minimal.

## Implementation Decisions

* **Storage Model:**
  * Each Group is stored as a flat YAML file (`<group_id>.yaml`) under the `requests/` directory.
  * A Group contains an optional `parent_id` string field referencing the ID of its parent Group. If empty, the Group is treated as a root-level entity.
  * Individual Requests are stored nested inline within their immediate parent Group's YAML file.
* **Tree Reconstruction:**
  * When loading requests, the system loads all Group files from disk and builds a tree structure in memory using the `parent_id` references.
  * Circular dependency detection is run on load. If a circular reference is found (e.g., Group A points to B, which points to A), loading fails with a clear error.
* **UI & Tree Updates:**
  * The GUI tree updates tree item labels based on whether they represent a Group or a Request.
  * Leaf nodes in the tree represent Requests; branch nodes represent Groups.

## Testing Decisions

* **Testing Philosophy:**
  * A good test verifies the structural integrity of the resolved tree hierarchy (including deep nesting, circular detection, and leaf mappings) by checking the parsed memory model rather than asserting GUI rendering.
* **Modules to Test:**
  * `internal/request` (specifically the tree loader and circular reference detector).
* **Prior Art:**
  * Unit tests in `internal/request/load_test.go`.

## Out of Scope

* Drag-and-drop tree reorganization in the GUI (reorganization will be handled via explicit context menus/dialogs).
* Linking requests to multiple groups simultaneously (a Request is owned strictly by a single Group).

## Further Notes

* A virtual `unsorted` Group will act as the catchment area for requests that do not belong to any user-defined group, and its requests will be displayed directly at the root of the tree.
