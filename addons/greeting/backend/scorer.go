// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"fmt"
	"strings"

	"github.com/margince/margince/backend/pkg/extension"
	"github.com/margince/margince/backend/pkg/extension/crm"
)

// scoreDeal is the deal scorer this unit provides: core has none of its own.
func scoreDeal(ctx context.Context, rt extension.Runtime, deal crm.Deal) (extension.DealScore, error) {
	var memberNotes int
	err := rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		return tx.QueryRow(ctx, `
			SELECT count(*) FROM ext.ext_acme_greeting_deal_note
			WHERE deal_id = $1::uuid AND kind = 'member'`, deal.Id).Scan(&memberNotes)
	})
	if err != nil {
		return extension.DealScore{}, err
	}
	return scoreFrom(deal, memberNotes), nil
}

// scoreFrom is Acme's rule: a deal is known when it names its company, has a
// value, and the account team has written about it.
func scoreFrom(deal crm.Deal, memberNotes int) extension.DealScore {
	score := 20
	var known []string
	if deal.CompanyId != nil {
		score += 30
		known = append(known, "names its company")
	}
	if deal.AmountMinor != nil {
		score += 30
		known = append(known, "has a value")
	}
	if counted := min(memberNotes, 4); counted > 0 {
		score += counted * 5
		known = append(known, fmt.Sprintf("has %d Acme note(s)", memberNotes))
	}
	if len(known) == 0 {
		return extension.DealScore{Score: score, Reason: "Acme knows nothing about this deal yet."}
	}
	return extension.DealScore{Score: score, Reason: "The deal " + strings.Join(known, ", ") + "."}
}
