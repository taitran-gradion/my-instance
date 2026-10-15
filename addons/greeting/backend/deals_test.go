// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"errors"
	"testing"

	"github.com/margince/margince/backend/pkg/extension"
	"github.com/margince/margince/backend/pkg/extension/crm"
)

// recordingDeals is core's deal service as the decorator sees it: it records
// the request it was handed.
type recordingDeals struct{ got *crm.CreateDealRequest }

func (r recordingDeals) CreateDeal(_ context.Context, in crm.CreateDealRequest) (crm.Deal, error) {
	*r.got = in
	return crm.Deal{Id: "01a13e33-cd5a-7058-8858-86739e83b2e6", Name: in.Name}, nil
}

// refusingTx is a runtime whose transaction fails, standing in for a welcome
// note that could not be written.
type refusingTx struct{ extension.Runtime }

var errNoteWrite = errors.New("the welcome note was not written")

func (refusingTx) Tx(context.Context, func(context.Context, extension.Tx) error) error {
	return errNoteWrite
}

func TestAPlaceholderNameIsRefusedBeforeCoreRuns(t *testing.T) {
	var got crm.CreateDealRequest
	_, err := decorateDeals(refusingTx{}, recordingDeals{&got}).CreateDeal(context.Background(), crm.CreateDealRequest{Name: " TBD "})
	if !errors.Is(err, extension.ErrInvalid) {
		t.Fatalf("err = %v, want ErrInvalid", err)
	}
	if got.Name != "" {
		t.Fatal("core's create ran for a refused deal")
	}
}

func TestALargeDealWithoutAPriorityReachesCoreAsHigh(t *testing.T) {
	var got crm.CreateDealRequest
	amount := int64(largeDealMinor)
	_, _ = decorateDeals(refusingTx{}, recordingDeals{&got}).CreateDeal(context.Background(), crm.CreateDealRequest{Name: "Fleet renewal", AmountMinor: &amount})
	if got.Priority == nil || *got.Priority != crm.CreateDealRequestPriorityHigh {
		t.Fatalf("priority = %v, want high", got.Priority)
	}
}

func TestAFailedWelcomeNoteFailsTheCreate(t *testing.T) {
	var got crm.CreateDealRequest
	_, err := decorateDeals(refusingTx{}, recordingDeals{&got}).CreateDeal(context.Background(), crm.CreateDealRequest{Name: "Fleet renewal"})
	if !errors.Is(err, errNoteWrite) {
		t.Fatalf("err = %v, want the note's failure, so core rolls the deal back", err)
	}
}
