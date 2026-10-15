// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"testing"

	"github.com/margince/margince/backend/pkg/extension/crm"
)

func TestScoreRisesWithWhatAcmeKnows(t *testing.T) {
	company, amount := "01a13e33-cd32-72ad-aaad-eb3410951dbb", int64(100)
	cases := []struct {
		name  string
		deal  crm.Deal
		notes int
		want  int
	}{
		{"nothing known", crm.Deal{}, 0, 20},
		{"company only", crm.Deal{CompanyId: &company}, 0, 50},
		{"company and value", crm.Deal{CompanyId: &company, AmountMinor: &amount}, 0, 80},
		{"notes count up to four", crm.Deal{CompanyId: &company, AmountMinor: &amount}, 9, 100},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := scoreFrom(tc.deal, tc.notes)
			if got.Score != tc.want || got.Reason == "" {
				t.Fatalf("score = %+v, want %d with a reason", got, tc.want)
			}
		})
	}
}
