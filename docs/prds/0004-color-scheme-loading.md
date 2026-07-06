# PRD-0004: Color Scheme Loading

## Problem Statement

I want to customize the look and feel of my REST client (TUI and GUI) to match my visual preferences and system aesthetics. Hardcoded light/dark colors can cause poor readability, strain my eyes, or clash with my operating system theme.

## Solution

Enable custom Color Scheme loading from the main `go-rest.yaml` configuration file. A scheme defines an abstract logical palette (e.g. background, foreground, primary, highlight, success, error) as HEX values. The application will map these abstract colors to the appropriate rendering properties (Lipgloss styling for TUI, and `fyne.Theme` values for GUI). Support pairing light/dark themes to allow dynamic, automatic switching based on system preferences.

## User Stories

1. As a user, I want to define my custom color schemes inline under the `color_schemes` section in the main `go-rest.yaml` file, so that I can manage my styles in a single configuration file.
2. As a user, I want to define a theme using a simple abstract palette of logical colors (background, foreground, primary, highlight, success, warning, error), so that I don't have to write separate design configurations for the terminal and the desktop GUI.
3. As a user, I want to select a single active `color_scheme` to force the application to use that theme statically.
4. As a user, I want to specify a `light_scheme` and a `dark_scheme` in my configuration, so that `go-rest` dynamically selects the corresponding theme based on whether my operating system is in Light or Dark mode.
5. As a developer, I want the theme engine to automatically map abstract colors to the correct framework variables (e.g., mapping `highlight` to Fyne's focus color or Lipgloss's border color).

## Implementation Decisions

* **Configuration Schema:**
  * Define `color_schemes` as a map in `go-rest.yaml` where each key is a theme name, containing:
    * `background`: primary background color
    * `foreground`: primary text color
    * `primary`: main accent/brand color
    * `highlight`: focus or selection highlight
    * `success`: success status color (e.g. HTTP 2xx)
    * `warning`: warning status color (e.g. HTTP 3xx)
    * `error`: error status color (e.g. HTTP 4xx/5xx)
  * Add configuration settings: `color_scheme`, `light_scheme`, and `dark_scheme`.
* **Dynamic Resolution:**
  * On launch and on system theme change notifications, the app checks if `light_scheme` and `dark_scheme` are configured.
  * If yes, it queries the host system theme and resolves to the matching scheme. If no system theme can be determined, or if `color_scheme` is set explicitly, it falls back to the static `color_scheme`.
* **Framework Integration:**
  * TUI: Convert Resolved active colors to `lipgloss.AdaptiveColor`.
  * GUI: Implement `fyne.Theme` interface in Go, wrapping the active scheme and returning color values dynamically based on palette lookups.

## Testing Decisions

* **Testing Philosophy:**
  * Test configuration parsing and active scheme resolution using Go unit tests. By passing mock YAML configurations and mocking the system theme status (light/dark mode flag), we can assert that the correct HEX values are returned by the theme resolver.
* **Modules to Test:**
  * `internal/config` (parsing) and `internal/ui/styles` (resolution and Fyne theme adapter).
* **Prior Art:**
  * Existing colors resolver in `internal/ui/styles/styles.go`.

## Out of Scope

* Importing external CSS or Fyne `.json` theme files directly.
* Real-time graphical theme designer inside the application (custom themes must be written directly in the YAML configuration file).
