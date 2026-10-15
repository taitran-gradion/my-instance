//go:build margince_overlay

// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

// Compiled into core's package deals, never into this extension's: the
// composer adds it with `go build -overlay`, and core's CreateDeal stays
// callable as CreateDeal__core.

package deals

import (
	"context"
	"time"

	crmcontracts "github.com/margince/margince/backend/internal/contracts"
)

// acmeCloseHorizon is when Acme expects a deal with no close date to close.
const acmeCloseHorizon = 30 * 24 * time.Hour

// CreateDeal gives a deal created without a close date Acme's usual one,
// marked provisional: a machine's proposal no human has confirmed, which the
// contract has no field for and only core's own create can set.
//
//margince:decorate
func (s *Store) CreateDeal(ctx context.Context, in CreateDealInput) (crmcontracts.Deal, error) {
	if in.ExpectedClose == nil {
		close := time.Now().UTC().Add(acmeCloseHorizon).Truncate(24 * time.Hour)
		in.ExpectedClose = &close
		in.CloseDateProvisional = true
	}
	return s.CreateDeal__core(ctx, in)
}
