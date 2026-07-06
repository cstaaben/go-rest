# CLI Subcommands and Ad-hoc Execution

We decided to support CLI-only operations to run requests without launching the GUI.
1. Running the bare binary with no arguments launches the GUI.
2. We support a minimalist set of subcommands: `run <request-identifier>` and `list`.
3. `run` supports flags (`--url`, `--method`, `--header`, `--body`) to construct and execute **Ad-hoc Requests** on-the-fly.
4. Ad-hoc requests are transient by default, but can be saved to the collection using the `--save` flag (with optional `--name` and `--group`).
5. Ad-hoc requests support **Variable Placeholders**, but execution fails immediately if any variable placeholder cannot be resolved (including if no environment is selected).
6. Output defaults to response body only, with flags `--include`/`-i`, `--verbose`/`-v`, and `--quiet`/`-q` (which suppresses all stdout output, e.g. for scripting where only exit status matters) to customize detail level.
