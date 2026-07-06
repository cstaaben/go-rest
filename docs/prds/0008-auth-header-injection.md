# PRD-0008: Authentication Header Injection

## Problem Statement

Most web APIs require some form of authentication (such as a Bearer Token, API Key, or Basic Auth). Currently, I must manually define these HTTP headers (e.g. writing `Authorization: Bearer my-token` or base64-encoding credentials for `Basic`) on every single request. This is highly repetitive, hard to maintain, and insecure since credentials get hardcoded directly in individual request files.

## Solution

Implement an Authentication Header Injection helper in the `internal/auth` package. Support three common auth configurations: **Basic Auth**, **Bearer Token**, and **API Key**. To eliminate duplicate configuration, implement **Authentication Inheritance** across our nested Group hierarchy, allowing a Request or sub-group to inherit its credentials from its parent Group. Ensure auth fields support variable placeholders to allow dynamic, safe credential swapping.

## User Stories

1. As a developer, I want to configure authentication for a Request or Group, so that the correct headers (or query parameters) are automatically constructed and injected when executing a request.
2. As a developer, I want to choose between standard authentication types: **Basic Auth** (username/password), **Bearer Token** (token string), and **API Key** (custom header or query parameter key/value), so that I can authenticate against most web APIs easily.
3. As a developer, I want my Requests (and sub-groups) to inherit authentication configurations from their parent Group by default, so that I can define my API credentials once at the collection root and share them across all endpoints.
4. As a developer, I want to explicitly override inherited credentials on any individual Request or child Group (including selecting a different auth type or choosing "No Auth"), so that I can manage public or mixed-authentication endpoints in the same collection.
5. As a developer, I want authentication fields (like tokens, API keys, or passwords) to support variable placeholders (e.g. `{{apiToken}}`), so that my secrets are resolved from my Active Environment at execution time and never saved directly to request files.

## Implementation Decisions

* **Auth Struct Model:**
  * Both `Request` and `Group` structs are updated to include an `Auth` field:
    * `type`: "basic" | "bearer" | "apikey" | "inherit" | "none" (default for Request is "inherit", default for root Group is "none").
    * `params`: A map of key-value parameters:
      * `basic`: `username`, `password`
      * `bearer`: `token`
      * `apikey`: `key`, `value`, `in` (where `in` is "header" or "query")
* **Inheritance Resolution Crawler:**
  * When executing a request, the engine traverses the Group tree using `parent_id` pointers to find the closest parent node whose `Auth.Type` is not `"inherit"`.
  * If the resolved configuration is `"none"`, no auth is applied. Otherwise, the resolved auth settings are used.
* **Header and Parameter Injection:**
  * Prior to sending the request, the resolved auth fields are evaluated through the Environment variable substitution engine.
  * The HTTP request is then decorated:
    * `basic`: Base64-encodes `username:password` and injects `Authorization: Basic <hash>`.
    * `bearer`: Injects `Authorization: Bearer <token>`.
    * `apikey`: Injects `key: value` into request headers, or appends `?key=value` to the request query string.

## Testing Decisions

* **Testing Philosophy:**
  * Verify the crawling inheritance chain and the generated header formatting using unit tests. Mock various nested Group hierarchies in memory and assert that:
    * The crawler resolves to the expected parent's auth config or "none".
    * The injected request headers contain the correct, valid HTTP structures (including base64 hashing verification).
    * Missing variables inside auth fields trigger the expected substitution failure.
* **Modules to Test:**
  * `internal/auth` (the crawler and injection formatter).
* **Prior Art:**
  * None.

## Out of Scope

* Dynamic token negotiation or signature algorithms (e.g., OAuth 2.0 authorization code flow, AWS Signature V4, digest authentication).
* Masking variable values in the CLI output (dealt with in the CLI PRD via quiet/verbose logs).
