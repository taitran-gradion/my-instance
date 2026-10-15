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

## The sample extension: `addons/greeting`

`acme-greeting` is laid out the way an extension is meant to be, and it shows every way an
extension changes Margince.

```
addons/greeting/
  go.mod                      the module; the replace in extensions.mod points here
  backend/                    the one Go package
    acmegreeting.go           New(): everything the extension gives core, declared in one place
    deals.go                  decorates core's deal service
    scorer.go                 provides the deal scorer core has no implementation of
    notes.go                  its own routes, over its own table
    hello.go                  the smallest route there is
    migrations/               its table, beside the code that embeds it
  api/crm.yaml                its routes, merged into core's contract
  frontend/
    index.tsx                 the package entry: its page, and `screens` it adds to core
    api.ts                    every call its screens make
    pages/                    #/ext/acme-greeting and #/x/acme-desk
    components/               the account modal
    overrides/                core files it replaces, pinned by overrides.lock
    i18n/                     its copy, in en, de and vi
```

| What | Where | What a user sees |
|---|---|---|
| Decorates a core service | `backend/deals.go` (`DecorateDeals`) | A deal named "TBD" is refused with Acme's reason; a deal of 50,000.00 or more gets priority High; every new deal gets a welcome note, written in the deal's own transaction |
| Provides a core service | `backend/scorer.go` (`ScoreDeals`) | `GET /deals/{id}/score` answers with Acme's score |
| Adds routes and a table | `backend/notes.go`, `api/crm.yaml`, `backend/migrations/` | Notes on a deal, each also logged on the deal's timeline, audited, and announced as `ext_acme_greeting.note_added` |
| Adds a page | `frontend/index.tsx` (`screens`) | "Acme desk" in the sidebar after Deals |
| Replaces a core file | `frontend/overrides/` | The "Acme account" badge on the deal header, opening the score and notes |

The notes routes are gated on `ext_acme_greeting_note`, which no role holds until an admin grants it:
Settings, Roles, the role, then `ext_acme_greeting_note` (or `PATCH /v1/roles/{key}/objects/ext_acme_greeting_note`).
Until then the modal shows the score and no notes.

After moving `go.mod` to a newer core, a failed `overrides.lock` line means: read core's new file,
update the override, then paste the printed line into the lock.
