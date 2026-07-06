# Authentication Header Injection

We decided to support automatic injection of HTTP headers and parameters for common authentication types:
1. We support three primary **Authentication Configurations**: **Basic Auth** (base64 encoded credentials in `Authorization: Basic`), **Bearer Token** (`Authorization: Bearer <token>`), and **API Key** (custom header or query parameter).
2. We support **Authentication Inheritance**: a Request or Group inherits the authentication configuration of its nearest parent Group by default. This can be explicitly overridden on any node by specifying a concrete auth configuration or "No Auth" (`none`).
3. Authentication fields support **Variable Placeholders** (resolved at execution time from the Active Environment).
