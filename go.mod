module example.com/acme-margince

go 1.27.2

require (
	example.com/acme-margince/addons/acme-greeting v0.0.0
	github.com/margince/margince/backend v0.0.0
	github.com/margince/margince/composition v0.0.0
	github.com/margince/margince/extensions/de v0.0.0
	github.com/margince/margince/extensions/openchannel v0.0.0
	github.com/margince/margince/extensions/vn v0.0.0
	github.com/taitran-gradion/my-extension v0.0.0-20261015042026-5239baace1de
)

// This installation's own extension, kept in this repository.
replace example.com/acme-margince/addons/acme-greeting => ./addons/greeting

// The composition this instance generated: core's own replace
// (backend/go.mod -> ../composition) is ignored outside the main module.
replace github.com/margince/margince/composition => ./.margince/core/build/composition/backend

// Core from the fork that carries the extensions.mod POC. The Makefile reads
// the backend version here to fetch the matching core for composing.
replace (
	github.com/margince/margince/backend => github.com/taitran-gradion/margince/backend v0.0.7-poc.2
	github.com/margince/margince/extensions/de => github.com/taitran-gradion/margince/extensions/de v0.0.7-poc.2
	github.com/margince/margince/extensions/openchannel => github.com/taitran-gradion/margince/extensions/openchannel v0.0.7-poc.2
	github.com/margince/margince/extensions/vn => github.com/taitran-gradion/margince/extensions/vn v0.0.7-poc.2
)
