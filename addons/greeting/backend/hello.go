// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"encoding/json"

	"github.com/margince/margince/backend/pkg/extension"
)

func hello(_ context.Context, _ extension.Runtime, _ json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(map[string]string{"message": "Hello from acme-greeting"})
}
