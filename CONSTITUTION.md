# Constitution: gcal-meet-969281

## Preamble: System Archetype
This is a lightweight, high-performance Go (Golang) CLI Utility designed to run locally or via scheduler (e.g., cron) to query Google Calendar and launch upcoming meetings. Reliability, secure token management, execution speed, and sanitization of external commands are paramount.

## Owned Domains
* Google Calendar API event integration, including OAuth2 consent, credential verification, and calendar event fetching.

* Meeting URL extraction from event descriptions, locations, or custom fields, along with target platform matching.

* Local OS execution and browser launching to open meeting links automatically.

## Forbidden Dependencies & Protected Zones
* NEVER import deprecated, unmaintained, or legacy third-party Google Calendar, API, or OAuth2 libraries.

* NEVER hardcode client secrets, OAuth2 tokens, or private credentials anywhere in the repository.

* DO NOT execute arbitrary, unsanitized shell commands or URLs directly to launch meetings.

## Security & Resilience Hardlines
* All target meeting URLs MUST be rigorously validated and sanitized to prevent command injection before OS browser/app execution.

* OAuth2 tokens and calendar secrets MUST be loaded securely from the environment or a secure local configuration file.

* Personal meeting details, user credentials, attendee emails, and tokens MUST NEVER be printed to standard output or logs.

## Architectural & Async Invariants
* All Google Calendar API interactions and network calls MUST enforce a strict timeout (e.g., 5 to 10 seconds).

* All network client operations MUST handle transient connectivity errors gracefully using exponential backoff or retry logic.

* All calendar fetching and API integrations MUST be decoupled from core matching and launching logic via clean Go interfaces.

## Observability & Error Bounding
* All errors MUST be wrapped with context using `%w` and propagated up to the CLI entry point.

* Structured logging via the standard library `slog` package MUST be used for tracking execution steps.

* Every logging statement MUST omit personal, identifying user details or private calendar event data in production configurations.

## Coding Conventions
* Follow standard Go styling (run `go fmt` and `go vet` before every commit).

* Use short, descriptive variable names in accordance with standard Go idiomatic patterns.

* Keep functions small, focused, and single-purpose, returning early on error conditions.

## Governance
* Constitution supersedes all other project documentation.

* Amendments require documentation and team approval.

* All PRs/reviews MUST verify compliance with these invariants.

**Version**: 1.0.0 | **Ratified**: 2026-06-29
