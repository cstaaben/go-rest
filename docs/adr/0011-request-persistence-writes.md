# Request Persistence Writes

We decided to support persistence operations on the filesystem as follows:
1. Modifying or deleting a Request overwrites the entire YAML file representing its parent Group (`<group_id>.yaml`).
2. Creating, renaming, or moving a Group overwrites its specific `<group_id>.yaml` metadata file.
3. Deleting a Group deletes its file on disk and recursively searches and deletes all nested child Group files pointing to it via `parent_id`.
