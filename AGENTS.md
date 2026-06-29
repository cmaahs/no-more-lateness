# Agent Governance: no-more-lateness

This document defines context routing, standard operating procedures, and strict meta-directives for AI agents collaborating on the `no-more-lateness` Go utility.

## Decision Protocol
Before making any architectural decision, you MUST state which section of CONSTITUTION.md governs that decision. If no section applies, flag it as uncharted area requiring team review.

## Context Routing
AI Agents must read the following files to align with repository guardrails and history before performing any implementation:
* **Invariants & Conventions**: Refer to [CONSTITUTION.md](./CONSTITUTION.md) for strict architectural, security, and coding invariants.
* **Architectural History**: Refer to [ARCHITECTURE.md](./ARCHITECTURE.md) to review past ADRs before proposing any changes or new abstractions.

## Agent SOPs

### 1. Jira / Ticket Tracking
* Include the Jira ticket ID in the branch name and commit messages (e.g., `feat/PROJ-123-add-widget`).

### 2. Git Branch & Commit Conventions
* **Branch Prefixes**: Use prefixes `feat/` (features), `fix/` (bug fixes), `chore/` (maintenance/tooling), or `hotfix/` (urgent production fixes).
* **Branching Strategy**: Always branch from `main` and open a PR/MR back to `main`.
* **Commit Messages**: Use conventional commit format: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`. Reference the Jira ticket ID in the commit message where applicable.

### 3. Cross-Cutting Process Rules
* All code changes MUST include corresponding tests (unit or integration) in a `_test.go` file within the same package.
* Before proposing changes, agents must verify that compiling, linting, and tests pass successfully.

### 4. Operational Commands
* **Build**: `go build -o no-more-lateness main.go`
* **Test**: `go test -v ./...`
* **Lint**: `golangci-lint run`

### Meta-Directives
* The Agent MUST NOT introduce new third-party dependencies unless explicitly mandated.
* The Agent MUST NOT implement features or speculative abstractions not explicitly requested.
* The Agent MUST NOT silently remove existing working code or tests to resolve errors.
* The Agent MUST append non-interactive and run-once flags (e.g., `--watch=false`, `CI=true`, `--no-interaction`, `-y`) to all test and build commands to prevent hanging in watch mode or on interactive prompts.
