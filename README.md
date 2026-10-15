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

`acme-greeting` is laid out the way an extension is meant to be, and it uses every level of
reach an extension has on Margince, each once.

```
addons/greeting/
  go.mod                      the module; the replace in extensions.mod points here
  backend/                    the one Go package
    acmegreeting.go           New(): everything the extension gives core, declared in one place
    hooks.go                  an operation hook and a write hook on new deals
    contacts.go               title-cased names, a task on every new contact, follow-ups
    scorer.go                 provides the deal scorer core has no implementation of
    notes.go                  its own routes, over its own table
    hello.go                  the smallest route there is
    migrations/               its table
    overlays/                 core functions it compiles over, pinned by overlays.lock
  api/crm.yaml                its routes, merged into core's contract
  frontend/
    index.tsx                 the package entry: its page, and `screens` it adds to core
    theme.css                 Acme's brand color over core's
    api.ts                    every call its screens make
    pages/                    #/ext/acme-greeting and #/x/acme-desk
    components/               the account modal
    overrides/                core files it replaces, pinned by overrides.lock
    i18n/                     its copy, in en, de and vi
```

| Level | Where | What a user sees |
|---|---|---|
| 1. Capability | `backend/scorer.go`, `backend/notes.go`, `frontend/index.tsx` | Acme's score at `GET /deals/{id}/score`; notes on a deal, each also on its timeline; "Acme desk" in the sidebar |
| 2. Write hook | `backend/hooks.go` (`welcomeNewDeal`) | Every new deal gets a welcome note, in the deal's own transaction: if the note fails, the deal is not created |
| 3. Operation hook | `backend/hooks.go` (`shapeNewDeal`) | A deal named "TBD" is refused with Acme's reason; a deal of 50,000.00 or more gets priority High |
| 4. Overlay | `backend/overlays/internal/modules/deals/` | A deal created with no close date gets one 30 days out, marked provisional. The contract has no field for that, so only core's own create can set it |
| Frontend override | `frontend/overrides/` | The "Acme account" badge on the deal header, opening the score and notes; on a contact's header, "Acme profile" (how much of email, phone, title and employer is known) and "Follow-ups" |

**Brand color.** `frontend/theme.css` redefines core's accent variables (`--accent`, `--accentText`,
`--accentLight`, `--accentMed`, `--glowA`) for light, dark and "dark when the system asks", with
selectors one step more specific than core's `design-system/tokens.css`. The extension's `index.tsx`
imports it, and core imports every extension's entry when the app starts, so the whole app is Acme
indigo instead of Margince emerald. Removing the import brings core's color back.

The contact page is customized with no change to core: an operation hook on `createContact` and
`updateContact` saves names title-cased, a write hook on `contact.create` files an "Introduce Acme" task
due in two days, and the follow-ups a member plans from the contact header are core tasks too, so core's
own "Needs attention" card shows them. The code is `backend/contacts.go`,
`backend/migrations/0002_contact_followup.*`, `frontend/overrides/src/screens/contactfacts.tsx` and
`frontend/components/followup-modal.tsx`.

Each level is listed in `manifest.generated.json` (`write_hooks`, `operation_hooks`, `overlays`,
`provides`). An overlay compiles code into core's own packages, so only an extension in this
repository may ship one; a published extension that does is refused at compose.

**Builds with overlays.** Go refuses to overlay a file in the module cache, so `make build` compiles
core from the clone in `.margince/core` (the same version `go.mod` names) whenever an extension ships
an overlay, and from the required module otherwise.

**Locks.** After moving `go.mod` to a newer core, a failed line in `backend/overlays.lock` or
`frontend/overrides.lock` means core changed what the extension was written against: read core's new
function or file, update the extension, then paste the printed line into the lock.

**Notes need a grant.** The notes routes are gated on `ext_acme_greeting_note`, which no role holds
until an admin grants it: Settings, Roles, the role, then read and create on `ext_acme_greeting_note`.
