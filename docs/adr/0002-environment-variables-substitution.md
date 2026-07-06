# Environment Variable Substitution

We decided that environment variables will be referenced using double curly braces (e.g. `{{variable_name}}`). During request execution:
1. Substitution is performed globally on the request URL, headers, and body.
2. If any placeholder cannot be resolved from the Active Environment, execution must fail immediately with a validation error identifying the missing variable.
3. The Active Environment is the sole source of truth; there is no fallback to OS/system environment variables.
