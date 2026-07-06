# go-rest Domain Context

A command-line and graphical REST client for sending HTTP requests and managing request collections.

## Language

**Request**:
A single HTTP request specification including method, URL, headers, protocol version, and optional body, along with its last cached response.
_Avoid_: REST call, query, endpoint, api call

**Request Group** (or **Group**):
A logical folder containing Requests and potentially other nested child Groups, forming a hierarchical tree.
_Avoid_: Collection, folder, request folder

**Unsorted Request**:
A Request that does not belong to any user-defined Group, which is stored in a special "unsorted" group but rendered directly at the root of the UI tree.
_Avoid_: Orphan request, raw request

**Environment**:
A named set of key-value variables used to substitute placeholders globally in URL, headers, or body when executing Requests.
_Avoid_: Config environment, workspace variables

**Active Environment**:
The single Environment currently selected by the user to resolve variables for request execution.
_Avoid_: Selected config, active variables

**Variable Placeholder**:
A reference to an Environment variable wrapped in double curly braces (e.g. `{{variable_name}}`) embedded within a Request's URL, headers, or body.
_Avoid_: Variable reference, dollar reference

**Saved Request**:
A Request stored permanently in a Group file on disk, accessible via the GUI tree or by name/UUID in the CLI.
_Avoid_: Stored request, database request

**Ad-hoc Request**:
A transient Request constructed on-the-fly using CLI command flags rather than being loaded from a saved collection.
_Avoid_: One-off request, transient request, CLI request

**Color Scheme**:
A named set of logical color definitions mapping abstract palette keys to HEX codes, used to style the user interface.
_Avoid_: Theme, skin, color palette

**Active Color Scheme**:
The Color Scheme currently applied to style the application, determined either by explicit selection or dynamically resolved based on system light/dark mode.
_Avoid_: Selected theme, active skin

**Request Snapshot**:
The fully-resolved state of a Request's URL, headers, and body at the moment of execution after variable substitution.
_Avoid_: Resolved request, sent request

**Request History**:
A persistent, rolling log of past executions of Saved Requests, detailing the exact Request Snapshot sent and response received.
_Avoid_: Log history, run history, run log

**Code Snippet**:
An equivalent code representation of a Request generated in a specific target language or library (e.g. `curl`, Go, Python, JavaScript) to allow execution outside of `go-rest`.
_Avoid_: Exported request, client code

**Snippet Target**:
A supported programming language or HTTP library format for generating Code Snippets.
_Avoid_: Language format, code type

**Import Spec**:
A schema definition file (such as an OpenAPI specification) containing API definitions that can be loaded into `go-rest` to generate Requests, Groups, and Environments.
_Avoid_: Schema file, api definition file

**Import Source**:
The local file path or remote URL of an Import Spec.
_Avoid_: Import file, import target

**Authentication Configuration** (or **Auth Config**):
The settings defined on a Request or Group to automatically inject credentials (headers or query parameters) prior to execution.
_Avoid_: Auth settings, login credentials

**Authentication Inheritance**:
The hierarchical resolution process where a Request or child Group automatically inherits the Authentication Configuration of its nearest parent Group.
_Avoid_: Parent auth, group inheritance

**Extraction Rule**:
A declarative instruction defined on a Request (under the **Variables** tab) to extract a specific value from its HTTP response (via JSONPath, XPath, or header name) and bind it to a Session Variable.
_Avoid_: Response query, data filter

**Session Variable**:
A temporary, in-memory variable populated dynamically during execution (such as by an Extraction Rule) that overrides active Environment variables during the current application session.
_Avoid_: Local variable, run variable, temp credential

**Response Data**:
The structured container (encapsulating status code, headers, body, protocol, and execution duration) representing the outcome of a Request's network execution.
_Avoid_: HTTP response, result data

**Dynamic Row List**:
A GUI key-value editor component consisting of a vertical stack of editable fields (Key, Value, Active checkbox, Delete button) with automatic trailing row generation.
_Avoid_: Parameter table, key-value grid

**Environment Manager**:
The modal GUI dialog used to create, rename, delete, and configure Environments and their variables.
_Avoid_: Config window, environment tab









