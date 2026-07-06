# GUI Parameter Editors

We decided to implement the following GUI layout components for Request editing:
1. HTTP headers, query parameters, and variable definitions are managed using a **Dynamic Row List**. This component displays key-value pairs vertically, with an active checkbox to temporarily toggle parameters, a delete button, and a trailing blank row that automatically spawns a new row upon input.
2. The request body is edited via a **Dynamic Body Selector** dropdown. Selecting `None` hides the editor; `JSON`, `Text`, or `XML` displays a scrollable text area; and `Form Data` displays a Dynamic Row List of form parameters.
