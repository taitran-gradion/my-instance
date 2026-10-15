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

// largeDealMinor is the amount, in minor units, from which Acme treats a new
// deal as high priority when its creator set none.
const largeDealMinor = 5_000_000

// placeholderNames are what Acme's team types before they know what a deal is.
var placeholderNames = []string{"tbd", "new deal", "deal", "test"}

// acmeDeals is core's deal service with Acme's rules around it.
type acmeDeals struct {
	rt   extension.Runtime
	next extension.DealService
}

func decorateDeals(rt extension.Runtime, next extension.DealService) extension.DealService {
	return acmeDeals{rt: rt, next: next}
}

// CreateDeal refuses a placeholder name, raises a large deal's priority, lets
// core create the deal, then files a welcome note in the same transaction.
func (d acmeDeals) CreateDeal(ctx context.Context, in crm.CreateDealRequest) (crm.Deal, error) {
	if isPlaceholder(in.Name) {
		return crm.Deal{}, fmt.Errorf("%w: Acme names a deal after what is being sold, not %q", extension.ErrInvalid, in.Name)
	}
	if in.Priority == nil && in.AmountMinor != nil && *in.AmountMinor >= largeDealMinor {
		high := crm.CreateDealRequestPriorityHigh
		in.Priority = &high
	}
	deal, err := d.next.CreateDeal(ctx, in)
	if err != nil {
		return crm.Deal{}, err
	}
	err = d.rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		_, err := insertNote(ctx, tx, deal.Id, "Welcome to Acme: an account manager reviews this deal within one working day.", "welcome", nil)
		return err
	})
	return deal, err
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
