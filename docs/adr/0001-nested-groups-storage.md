# Nested Groups Storage Model

We decided to support nested Request Groups by using flat Group files (`<group_id>.yaml`) under the `/requests/` directory. Each Group file can contain an optional `parent_id` field referencing its parent Group, and we reconstruct the hierarchy in memory. Requests remain nested inside their parent Group's YAML file.

We chose this model over nested directories or inline nested YAML trees to avoid managing complex directory hierarchies and to keep file sizes and change diffs small.
