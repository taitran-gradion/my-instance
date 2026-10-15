// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/margince/margince/backend/pkg/extension"
)

func shape(t *testing.T, body string) (*extension.Operation, error) {
	t.Helper()
	op := &extension.Operation{ID: "createDeal", Body: []byte(body)}
	return op, shapeNewDeal(context.Background(), nil, op)
}

func TestAPlaceholderNameIsRefused(t *testing.T) {
	if _, err := shape(t, `{"name":" TBD "}`); !errors.Is(err, extension.ErrInvalid) {
		t.Fatalf("err = %v, want ErrInvalid", err)
	}
}

func TestALargeDealWithoutAPriorityReachesCoreAsHigh(t *testing.T) {
	op, err := shape(t, `{"name":"Fleet renewal","amount_minor":5000000}`)
	if err != nil || !strings.Contains(string(op.Body), `"priority":"high"`) {
		t.Fatalf("body %s, err %v; want priority high", op.Body, err)
	}
}

func TestAPriorityTheCreatorSetIsKept(t *testing.T) {
	op, err := shape(t, `{"name":"Fleet renewal","amount_minor":9000000,"priority":"low"}`)
	if err != nil || !strings.Contains(string(op.Body), `"priority":"low"`) {
		t.Fatalf("body %s, err %v; want the creator's priority", op.Body, err)
	}
}

// refusingTx stands in for a welcome note that could not be written.
type refusingTx struct{ extension.Runtime }

var errNoteWrite = errors.New("the welcome note was not written")

func (refusingTx) Tx(context.Context, func(context.Context, extension.Tx) error) error {
	return errNoteWrite
}

func TestAFailedWelcomeNoteFailsTheWrite(t *testing.T) {
	err := welcomeNewDeal(context.Background(), refusingTx{}, extension.Write{Entity: "deal", Action: "create", ID: "01a13e33-cd5a-7058-8858-86739e83b2e6"})
	if !errors.Is(err, errNoteWrite) {
		t.Fatalf("err = %v, want the note's failure, so core rolls the deal back", err)
	}
}
