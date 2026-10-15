// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

// Package acmegreeting is this installation's own extension, kept in the
// instance repository and turned on by a local replace in extensions.mod.
//
// It shows each way an extension changes Margince:
//
//	deals.go      decorates core's deal service: refuses a placeholder name and
//	              files a welcome note, in the deal's own transaction
//	scorer.go     provides the deal scorer core has no implementation of
//	notes.go      its own routes, over its own table, with a ledger entry and a
//	              timeline activity on the deal for each note
//	hello.go      the smallest route there is
//	migrations/   its table, beside the code that embeds it
//	../api/       its routes, merged into core's contract
//	../frontend/  its pages and the deal header it decorates
package acmegreeting

import (
	"embed"

	"github.com/margince/margince/backend/pkg/extension"
)

//go:embed migrations
var migrations embed.FS

// New returns the unit's declaration: everything it gives core, in one place.
func New() extension.Extension {
	return extension.Extension{
		Name:        "acme-greeting",
		Version:     "0.2.0",
		Description: "This installation's own extension: account notes on deals, a deal score, and a decorated deal create.",
		Tools: []extension.Tool{
			{Name: "acme_greeting_hello", Handle: hello},
			{Name: "acme_greeting_list_notes", Handle: listNotes},
			{Name: "acme_greeting_add_note", Handle: addNote},
		},
		// Lists the unit under Settings, Integrations.
		Secrets: []extension.SecretsRequest{
			{Key: "greeting_token", Scope: extension.SecretScopeWorkspace},
		},
		DecorateDeals: decorateDeals,
		ScoreDeals:    scoreDeal,
		Migrations:    migrations,
	}
}
