# Response Extraction and Request Chaining

We decided to support extracting values from responses to enable request chaining:
1. Extraction rules are defined declaratively on a Request. In the GUI, this is configured via a **Variables** tab in the Request Editor. In the CLI, they can be specified using `--extract` flags (e.g. during ad-hoc run or saved request modification).
2. We support three query types:
   * **JSONPath** for JSON response bodies.
   * **XPath** for XML response bodies. To eliminate XML namespace complexity for the user, the engine automatically compiles simple slash-delimited paths (e.g. `/Envelope/Body/Token`) into namespace-agnostic wildcards (e.g. `/*[local-name()='Envelope']/*[local-name()='Body']/*[local-name()='Token']`) before evaluation.
   * **Header name lookup** for response headers.
3. Extracted values are saved as **Session Variables** in a temporary, in-memory store.
4. Session Variables override active Environment variables during the execution session and are discarded when the application exits.
