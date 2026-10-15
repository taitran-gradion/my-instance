# margince-instance

A test Margince installation that requires core as a Go module instead of
cloning it (margince/margince#7727, POC).

| File | What it holds |
|---|---|
| `go.mod` | requires core's `backend` and links this instance's own composition |
| `extensions.mod`, `extensions.sum` | the extensions this instance turns on, pinned by version and checksum |
| `config/` | this instance's deployment config |

`make up` runs everything: compose, build, database, migrate, api, worker, web.
`make stop` stops it. `make help` lists each step.

Core is the fork `github.com/taitran-gradion/margince`, which carries the
`extensions.mod` POC. `go.mod` replaces core's modules with the fork's at
`v0.0.7-poc.1` (tags `backend/v0.0.7-poc.1`, `extensions/<name>/v0.0.7-poc.1`).
To move core, change that version in `go.mod`; the Makefile reads it from there.

Stand-ins until core publishes the `margince` CLI and the SPA package: `make core`
clones the fork at the same tag into `.margince/core` (ignored), and composing and
the web app run from that clone.
