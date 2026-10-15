// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

// Package acmegreeting is this installation's own extension: it lives in the
// instance repository, is never published, and is turned on by a local
// replace in extensions.mod. It serves one greeting and a screen.
package acmegreeting

import (
	"context"
	"encoding/json"

	"github.com/margince/margince/backend/pkg/extension"
)

// New returns the unit's declaration.
func New() extension.Extension {
	return extension.Extension{
		Name:        "acme-greeting",
		Version:     "0.1.0",
		Description: "This installation's own extension: a greeting route and a screen, kept in the instance repository.",
		Tools: []extension.Tool{
			{Name: "acme_greeting_hello", Handle: hello},
		},
		// Lists the unit under Settings, Integrations, beside my-extension.
		Secrets: []extension.SecretsRequest{
			{Key: "greeting_token", Scope: extension.SecretScopeWorkspace},
		},
	}
}

func hello(_ context.Context, _ extension.Runtime, _ json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(map[string]string{"message": "Hello from acme-greeting"})
}
