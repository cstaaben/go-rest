# PRD-0012: GUI Parameter Editors

## Problem Statement

I cannot interactively configure request parameters, headers, or request bodies in the GUI. The editor tabs are static placeholder text labels, forcing me to close the app and manually edit the YAML files on disk to configure headers or payloads.

## Solution

Implement interactive GUI parameter editors in the Fyne desktop client. Create a reusable **Dynamic Row List** container for editing headers, query parameters, cookies, and variables. Create a **Dynamic Body Selector** that provides a dropdown to select the payload format (`None`, `JSON`, `Text`, `XML`, `Form Data`) and displays the correct entry layout dynamically.

## User Stories

1. As a user, I want to add, edit, and delete request headers and query parameters in a grid-like key-value list in the GUI, so that I don't have to manually format YAML text.
2. As a user, I want each key-value parameter row to feature an "Active" checkbox, so that I can temporarily disable headers (like authorization or cache-control headers) without deleting them.
3. As a user, I want a new empty row to be automatically appended to the bottom of the list when I begin typing in the last empty row, so that adding multiple headers is fluid and requires no extra clicks.
4. As a user, I want a delete button next to each key-value row to remove it instantly, so that I can easily clean up unused parameters.
5. As a user, I want to choose my request body type from a dropdown menu (None, JSON, Text, XML, Form Data), so that the input interface adapts to the payload format.
6. As a user, I want selecting `None` as the body type to hide the body input field completely, keeping my workspace clean for GET requests.
7. As a user, I want selecting `JSON`, `Text`, or `XML` to display a scrollable multiline text input, so that I can write structured payloads.
8. As a user, I want selecting `Form Data` to display a Dynamic Row List of parameters, so that I can build `multipart/form-data` request bodies.

## Implementation Decisions

* **Dynamic Row List Widget:**
  * Implement the key-value list as a custom vertical box container (`container.NewVBox`) dynamically populated with horizontal row containers (`container.NewHBox`).
  * Each row displays:
    * `Active`: `widget.NewCheck`
    * `Key`: `widget.NewEntry`
    * `Value`: `widget.NewEntry`
    * `Delete`: `widget.NewButton` (with a trash icon)
  * Bind an `OnChanged` listener on the Key/Value fields of the last row. If either gets text, append a new empty row to the parent VBox and refresh the layout.
* **Dynamic Body Selector:**
  * Implement the Body tab layout using a vertical box container.
  * At the top: A `widget.Select` dropdown containing `["None", "JSON", "Text", "XML", "Form Data"]`.
  * Below: An adaptive container holding the inputs. Swapping the dropdown selection hides/shows the appropriate input widget (a `widget.NewMultiLineEntry` or a Dynamic Row List VBox).

## Testing Decisions

* **Testing Philosophy:**
  * Test the widget controllers and the request-data mapper functions using Fyne's built-in testing framework (`fyne.io/fyne/v2/test`). Assert that:
    * Adding characters to the last row appends a new empty row to the list.
    * Toggling the active checkbox or clicking delete modifies the underlying `request.Data` struct correctly.
    * Changing the body dropdown selection accurately swaps visible panels and modifies the outgoing request content-type header.
* **Modules to Test:**
  * `internal/model/container.go` (and custom widget panels).
* **Prior Art:**
  * None.

## Out of Scope

* Advanced editor enhancements like JSON syntax highlighting, auto-completion, or bracket matching inside the body text area (plain text multiline input is sufficient for the prototype).
* Binary file upload mappings in `Form Data` (only text parameters are supported in this version).
