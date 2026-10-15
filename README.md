# margince-instance

A test Margince installation that requires core as a Go module instead of
cloning it (margince/margince#7727, POC).

| File | What it holds |
|---|---|
| `go.mod` | requires core's `backend` and links this instance's own composition |
| `extensions.mod`, `extensions.sum` | the extensions this instance turns on, pinned by version and checksum |
| `config/` | this instance's deployment config |
| `addons/greeting/` | `acme-greeting`, this installation's own extension, never published |

`make up` runs everything: compose, build, database, migrate, api, worker, web.
`make stop` stops it. `make help` lists each step.

Core is the fork `github.com/taitran-gradion/margince`, which carries the
`extensions.mod` POC. `go.mod` replaces core's modules with the fork's at
`v0.0.7-poc.1` (tags `backend/v0.0.7-poc.1`, `extensions/<name>/v0.0.7-poc.1`).
To move core, change that version in `go.mod`; the Makefile reads it from there.

Stand-ins until core publishes the `margince` CLI and the SPA package: `make core`
clones the fork at the same tag into `.margince/core` (ignored), and composing and
the web app run from that clone.

## Extensions

`extensions.mod` lists every extension this instance turns on, wherever its code lives:

```text
require (
	example.com/acme-margince/addons/acme-greeting v0.0.0                        // ours, in addons/greeting
	github.com/taitran-gradion/my-extension v0.0.0-20261015042026-5239baace1de   // public, from GitHub
)
replace example.com/acme-margince/addons/acme-greeting => ./addons/greeting
```

A public extension is downloaded and checked against `extensions.sum`. A `replace` onto a folder
is composed where it sits, so an edit under `addons/greeting/frontend/` reloads in the browser and
a Go change needs `make up` again. No folder name is fixed: the `replace` says where the code is,
and the extension's name comes from its module path (`acme-greeting`), not from its folder.
`go.mod` carries the same `require` and `replace`, because Go builds the binaries from it.
