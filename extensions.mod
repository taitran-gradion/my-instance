module margince.local/installation

go 1.27.2

require (
	example.com/acme-margince/addons/acme-greeting v0.0.0
	github.com/taitran-gradion/my-extension v0.0.0-20261015042026-5239baace1de
)

// This installation's own extension: composed in place, never downloaded.
replace example.com/acme-margince/addons/acme-greeting => ./addons/greeting
