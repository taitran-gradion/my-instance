module example.com/acme-margince

go 1.27.2

require (
	example.com/margince-ext/my-extension v0.3.1
	github.com/margince/margince/backend v0.0.0
	github.com/margince/margince/composition v0.0.0
	github.com/margince/margince/extensions/de v0.0.0
	github.com/margince/margince/extensions/openchannel v0.0.0
	github.com/margince/margince/extensions/vn v0.0.0
)

// The composition this instance generated: core's own replace
// (backend/go.mod -> ../composition) is ignored outside the main module.
replace github.com/margince/margince/composition => ./build/composition/backend

// STAND-INS until core publishes versions, all pointing at the core toolkit
// checkout named in core.path: backend for a backend/vX tag, and the units
// core bundles for their own module versions.
replace (
	github.com/margince/margince/backend => ../margince-poc-gomod/backend
	github.com/margince/margince/extensions/de => ../margince-poc-gomod/extensions/de
	github.com/margince/margince/extensions/openchannel => ../margince-poc-gomod/extensions/openchannel
	github.com/margince/margince/extensions/vn => ../margince-poc-gomod/extensions/vn
)
