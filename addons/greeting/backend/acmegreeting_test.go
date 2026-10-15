// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNewNamesTheUnitAfterItsModulePath(t *testing.T) {
	if got := New().Name; got != "acme-greeting" {
		t.Fatalf("Name = %q, want acme-greeting: the composer takes the name from the module path", got)
	}
}

func TestHelloGreets(t *testing.T) {
	out, err := hello(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("hello: %v", err)
	}
	var body struct{ Message string }
	if err := json.Unmarshal(out, &body); err != nil || body.Message == "" {
		t.Fatalf("hello returned %s, want a message", out)
	}
}
