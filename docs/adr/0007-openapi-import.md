# OpenAPI Specification Import

We decided to support importing API definitions from OpenAPI specifications:
1. In the GUI, the feature is accessed via the `File > Import...` menu item, which prompts the user for an **Import Source** (local file path or remote URL).
2. In the CLI, the feature is triggered using the `go-rest import <path-or-url>` subcommand.
3. The import maps the OpenAPI API title to a parent Group. Endpoints with `tags` are organized into nested child Groups named after each tag under the parent.
4. OpenAPI `servers` URLs are converted into separate Environment files containing a `baseUrl` variable.
5. Path parameters (e.g., `{id}`) in endpoint URLs are converted to Variable Placeholders (e.g., `{{id}}`).
