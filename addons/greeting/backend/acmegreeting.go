// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

// Package acmegreeting is this installation's own extension, kept in the
// instance repository and turned on by a local replace in extensions.mod.
//
// It shows each way an extension changes Margince:
//
//	hooks.go      an operation hook before core's createDeal, and a write hook
//	              inside the transaction that creates a deal
//	contacts.go   title-cased contact names, an "Introduce Acme" task on every
//	              new contact, and follow-ups planned from the contact page
//	scorer.go     provides the deal scorer core has no implementation of
//	notes.go      its own routes, over its own table, with a ledger entry and a
//	              timeline activity on the deal for each note
//	hello.go      the smallest route there is
//	migrations/   its table, beside the code that embeds it
//	../api/       its routes, merged into core's contract
//	overlays/     core functions it compiles over (this installation only)
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
		Version:     "0.4.0",
		Description: "This installation's own extension: account notes on deals, a deal score, and Acme's rules on new deals.",
		Tools: []extension.Tool{
			{Name: "acme_greeting_hello", Handle: hello},
			{Name: "acme_greeting_list_notes", Handle: listNotes},
			{Name: "acme_greeting_add_note", Handle: addNote},
			{Name: "acme_greeting_list_followups", Handle: listFollowups},
			{Name: "acme_greeting_plan_followup", Handle: planFollowup},
		},
		// Lists the unit under Settings, Integrations.
		Secrets: []extension.SecretsRequest{
			{Key: "greeting_token", Scope: extension.SecretScopeWorkspace},
		},
		OperationHooks: []extension.OperationHook{
			{Operations: []string{"createDeal"}, Handle: shapeNewDeal},
			{Operations: []string{"createContact", "updateContact"}, Handle: titleCaseNames},
		},
		WriteHooks: []extension.WriteHook{
			{Entity: "deal", Actions: []string{"create"}, Handle: welcomeNewDeal},
			{Entity: "contact", Actions: []string{"create"}, Handle: introduceAcme},
		},
		ScoreDeals: scoreDeal,
		Migrations: migrations,
	}
}
