# Code Snippet Generation

We decided to support generating equivalent programming code snippets from Request definitions:
1. Snippets can be generated in a set of core **Snippet Targets**: `curl`, Go (`net/http`), Python (`requests`), and JavaScript (`fetch`).
2. Snippet generation supports a toggle between generating with **Variable Placeholders** intact (Raw) or fully resolved with active environment variables (Resolved), defaulting to Resolved.
3. In the GUI, the generator is triggered by a button in the header next to the Send button, opening a modal containing the code, language selector, copy button, and raw/resolved toggle.
4. In the CLI, the generator is invoked via the `go-rest snippet <request-identifier>` subcommand, accepting `--lang`/`-l` and `--raw` flags.
