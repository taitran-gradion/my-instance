// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"strings"
	"testing"

	"github.com/margince/margince/backend/pkg/extension"
)

func TestNamesAreSavedTitleCased(t *testing.T) {
	op := &extension.Operation{ID: "updateContact", Body: []byte(`{"full_name":"trab van tran","first_name":"trab","last_name":"van tran","title":"cto"}`)}
	if err := titleCaseNames(context.Background(), nil, op); err != nil {
		t.Fatal(err)
	}
	body := string(op.Body)
	if !strings.Contains(body, `"full_name":"Trab Van Tran"`) || !strings.Contains(body, `"first_name":"Trab"`) || !strings.Contains(body, `"last_name":"Van Tran"`) || !strings.Contains(body, `"title":"cto"`) {
		t.Fatalf("body %s: want both names title-cased and nothing else touched", body)
	}
}

func TestAnInnerCapitalIsKept(t *testing.T) {
	if got := titleCase("mcDonald"); got != "McDonald" {
		t.Fatalf("titleCase = %q, want McDonald", got)
	}
}

func TestABodyWithNothingToChangeIsLeftAsSent(t *testing.T) {
	sent := `{"first_name": "Trab"}`
	op := &extension.Operation{ID: "updateContact", Body: []byte(sent)}
	if err := titleCaseNames(context.Background(), nil, op); err != nil || string(op.Body) != sent {
		t.Fatalf("body %s, err %v; want it untouched", op.Body, err)
	}
}

// callerless is a runtime with nobody behind it, as an import job runs.
type callerless struct{ extension.Runtime }

func (callerless) Caller() extension.Caller { return extension.Caller{} }

func TestAContactImportedWithNobodyBehindItGetsNoTask(t *testing.T) {
	if err := introduceAcme(context.Background(), callerless{}, extension.Write{Entity: "contact", Action: "create", ID: "01a13ea3-4232-765b-a65b-6b4e19671660"}); err != nil {
		t.Fatalf("err = %v, want nil: an import must not fail on Acme's task", err)
	}
}
