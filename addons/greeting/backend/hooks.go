// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/margince/margince/backend/pkg/extension"
)

// largeDealMinor is the amount, in minor units, from which Acme treats a new
// deal as high priority when its creator set none.
const largeDealMinor = 5_000_000

// placeholderNames are what Acme's team types before they know what a deal is.
var placeholderNames = []string{"tbd", "new deal", "deal", "test"}

// shapeNewDeal runs before core's createDeal: it refuses a placeholder name
// and raises a large deal's priority, by changing the body core then reads.
func shapeNewDeal(_ context.Context, _ extension.Runtime, op *extension.Operation) error {
	var body map[string]any
	if err := json.Unmarshal(op.Body, &body); err != nil {
		return fmt.Errorf("%w: the deal is not a JSON object", extension.ErrInvalid)
	}
	if name, _ := body["name"].(string); isPlaceholder(name) {
		return fmt.Errorf("%w: Acme names a deal after what is being sold, not %q", extension.ErrInvalid, name)
	}
	amount, hasAmount := body["amount_minor"].(float64)
	if _, hasPriority := body["priority"]; hasPriority || !hasAmount || amount < largeDealMinor {
		return nil
	}
	body["priority"] = "high"
	changed, err := json.Marshal(body)
	if err != nil {
		return err
	}
	op.Body = changed
	return nil
}

// welcomeNewDeal runs inside the transaction that created the deal, so the
// deal and its welcome note commit together or not at all.
func welcomeNewDeal(ctx context.Context, rt extension.Runtime, w extension.Write) error {
	return rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		_, err := insertNote(ctx, tx, w.ID, "Welcome to Acme: an account manager reviews this deal within one working day.", "welcome", nil)
		return err
	})
}

func isPlaceholder(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, placeholder := range placeholderNames {
		if name == placeholder {
			return true
		}
	}
	return false
}
