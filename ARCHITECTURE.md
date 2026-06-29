# Architecture Decision Log (ADRs)

This document captures past and present architectural choices and rejected designs for the `no-more-lateness` Go utility.

## ADR Index

* [ADR-001: Adoption of Agentic SDLC](#adr-001-adoption-of-agentic-sdlc)

---

## ADR-001: Adoption of Agentic SDLC

### Status
Accepted

### Context
The `no-more-lateness` repository is a newly initialized Go (Golang) CLI utility with a highly focused scope: querying Google Calendar and spawning local browser processes to launch upcoming online meetings. To ensure that AI-driven development agents maintain absolute alignment with critical requirements (such as the absolute prevention of credential leakage, sanitization of local OS commands, and clean interfaces for API wrappers), we need a rigid governance and configuration system that directs agent behavior and establishes invariants.

### Decision
We adopt the Alteryx Agentic SDLC Standard for repository governance. This introduces three golden governance files:
1. `GEMINI.md`: Bootloader pointing directly to `AGENTS.md`.
2. `CONSTITUTION.md`: Defines strict, non-negotiable invariants (such as secure token loading, URL sanitization, and interface-driven API separation) alongside general team conventions.
3. `AGENTS.md`: Outlines agent operational rules, conventional workflow patterns (Jira key integration in git), and mandates the strict validation of all architectural decisions against `CONSTITUTION.md`.

All development agents interacting with this codebase MUST read and strictly adhere to these governance files.

### Consequences
* **Aesthetic and Functional Safety**: Agents are strictly forbidden from executing un-sanitized CLI launcher commands or writing legacy/un-vetted OAuth library integrations.
* **Deterministic Operations**: Agents will always execute non-interactive build, lint, and test commands (`go build`, `go test -v ./...`, and `golangci-lint run`) avoiding watch-mode hang-ups.
* **Traceable Architecture**: All modifications proposed by agents must cite the specific section of the `CONSTITUTION.md` that governs the change, or flag it as uncharted territory requiring human-in-the-loop review.
