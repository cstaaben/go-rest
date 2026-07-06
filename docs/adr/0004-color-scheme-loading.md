# Color Scheme Loading

We decided to support custom user-defined Color Schemes:
1. Schemes are defined inline in the main configuration file (`go-rest.yaml`) under `color_schemes`.
2. Each scheme is a flat, single-mode palette mapping abstract logical colors (`background`, `foreground`, `primary`, `highlight`, `success`, `warning`, `error`) to HEX values.
3. The configuration supports setting a static `color_scheme`, or pairing `light_scheme` and `dark_scheme` together to dynamically apply themes based on the host system's light/dark mode.
4. TUI and GUI rendering modules map these abstract logical colors to Lipgloss style properties and Fyne theme variables respectively.
