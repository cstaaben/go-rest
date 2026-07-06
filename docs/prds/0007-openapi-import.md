# PRD-0007: OpenAPI Specification Import

## Problem Statement

When I want to test an existing API that is already fully defined in an OpenAPI (Swagger) specification, I have to manually recreate every single request, header, group, and server variable inside my REST client. This is extremely slow, tedious, and prone to transcription errors, and it gets out of sync whenever the API schema changes.

## Solution

Implement an OpenAPI Spec Importer in the `internal/importer` package. The importer will parse an OpenAPI 3.0/Swagger schema from a local file or remote URL, map the API definition to a parent Group, organize requests into nested sub-groups based on their OpenAPI `tags`, rewrite path parameters to `{{placeholder}}` syntax, and convert the spec's `servers` list into working Environment files. Access this via a `File > Import...` menu in the GUI and a `go-rest import` subcommand in the CLI.

## User Stories

1. As a developer, I want to import API definitions from an OpenAPI specification (JSON or YAML) from a local file path or remote URL, so that I don't have to manually recreate my requests.
2. As a developer, I want the importer to create a parent Group named after the API title, and place endpoints into nested child Groups based on their OpenAPI `tags`, so that my sidebar tree matches the structure of the API documentation.
3. As a developer, I want the importer to convert the OpenAPI `servers` section into separate `go-rest` Environment files (e.g. `<API Title> Production` and `<API Title> Staging`) containing a `baseUrl` variable, so that I have swappable endpoint configurations right out of the box.
4. As a developer, I want the importer to translate path parameters (e.g. `{userId}`) in endpoint URLs into `go-rest` Variable Placeholders (e.g., `{{userId}}`), so that they hook into our variable resolution and validation system.
5. As a developer, I want to access the import tool in the GUI via `File > Import...`, which displays a dialog allowing me to browse local files or input a schema URL.
6. As a developer, I want to run `go-rest import <path-or-url>` from the CLI, so that I can automatically script imports or set up git hooks to update my local request collection when the server schema updates.

## Implementation Decisions

* **Importer Module:**
  * Implement the parsing and mapping logic in the `internal/importer` package.
  * Use a Go YAML/JSON parser to read the OpenAPI document structure.
  * Map paths and method operations to `Request` objects, converting request bodies, query strings, and headers to equivalent request structures.
* **Hierarchy and Variable Re-writing:**
  * Tag Mapping: If a path operation has tags, create a child `Group` file for each unique tag (with `parent_id` referencing the parent API Group) and place the request inside it.
  * Path Param Regex: Search path strings for `{var}` patterns and replace them with `{{var}}`.
  * Server Map: Parse `servers` entries. For each entry, create an Environment struct with name `[API Title] - [Server Description/URL]` and write a key `baseUrl` with the server URL.
* **GUI Menu Action:**
  * Add `File > Import...` to the Fyne menu bar. Clicking it opens a generic dialog with a file picker and URL text input.

## Testing Decisions

* **Testing Philosophy:**
  * Test the importer logic by writing Go unit tests that pass mock OpenAPI spec payloads (JSON and YAML string buffers) to the importer. Assert that:
    * The resolved Group slices match the expected tag-based hierarchy.
    * URLs are correctly translated to use `{{baseUrl}}`.
    * Path variables are rewritten to use `{{placeholder}}`.
    * Output environment files match the server specifications.
* **Modules to Test:**
  * `internal/importer` (the schema compiler).
* **Prior Art:**
  * None.

## Out of Scope

* Importing other client collections (e.g., Postman collections, Insomnia exports) in this initial PRD version (though the GUI menu title "Import..." is kept generic to accommodate future import sources).
* Full parsing of OpenAPI request body schemas to auto-generate complex JSON payloads (generate empty body structure with basic types where possible, but deep mock payload generation is out of scope).
