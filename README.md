# margince-instance

A test Margince installation that requires core as a Go module instead of
cloning it (margince/margince#7727, POC).

| File | What it holds |
|---|---|
| `go.mod` | requires core's `backend` and links this instance's own composition |
| `extensions.mod`, `extensions.sum` | the extensions this instance turns on, pinned by version and checksum |
| `config/` | this instance's deployment config |
| `core.path` | the core checkout used for what core does not publish yet |

`make up` runs everything: compose, build, database, migrate, api, worker, web.
`make stop` stops it. `make help` lists each step.

Stand-ins until core publishes versions: the `backend` version and the
extensions core bundles are `replace` lines to the checkout in `core.path`,
composing runs that checkout's `make composition`, and the web app is served by
its Vite.
