// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/margince/margince/backend/pkg/extension"
	"github.com/margince/margince/backend/pkg/extension/crm"
)

const (
	noteEntity    = "ext_acme_greeting_deal_note"
	maxNoteLength = 2000
	recentNotes   = 20
)

type note struct {
	ID        string    `json:"id"`
	DealID    string    `json:"deal_id"`
	Body      string    `json:"body"`
	Kind      string    `json:"kind"`
	AuthorID  *string   `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// listNotes answers GET /ext/acme-greeting/notes: one deal's notes, or the
// most recent across deals when no deal_id is sent.
func listNotes(ctx context.Context, rt extension.Runtime, in json.RawMessage) (json.RawMessage, error) {
	args, err := extension.DecodeArgs[struct {
		DealID *string `json:"deal_id"`
	}](in)
	if err != nil {
		return nil, err
	}
	if args.DealID != nil && !extension.IsCanonicalUUID(*args.DealID) {
		return nil, fmt.Errorf("%w: deal_id is not a deal id", extension.ErrInvalid)
	}
	notes := []note{}
	err = rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id::text, deal_id::text, body, kind, author_id::text, created_at
			FROM ext.ext_acme_greeting_deal_note
			WHERE $1::uuid IS NULL OR deal_id = $1::uuid
			ORDER BY created_at DESC
			LIMIT $2`, args.DealID, recentNotes)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var n note
			if err := rows.Scan(&n.ID, &n.DealID, &n.Body, &n.Kind, &n.AuthorID, &n.CreatedAt); err != nil {
				return err
			}
			notes = append(notes, n)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string][]note{"notes": notes})
}

// addNote answers POST /ext/acme-greeting/notes. The note also lands on the
// deal's timeline as a core activity, which is what checks that the caller may
// see the deal: a deal they cannot see refuses the whole write.
func addNote(ctx context.Context, rt extension.Runtime, in json.RawMessage) (json.RawMessage, error) {
	args, err := extension.DecodeArgs[struct {
		DealID string `json:"deal_id"`
		Body   string `json:"body"`
	}](in)
	if err != nil {
		return nil, err
	}
	body := strings.TrimSpace(args.Body)
	switch {
	case !extension.IsCanonicalUUID(args.DealID):
		return nil, fmt.Errorf("%w: deal_id is not a deal id", extension.ErrInvalid)
	case body == "" || len(body) > maxNoteLength:
		return nil, fmt.Errorf("%w: a note is 1 to %d characters", extension.ErrInvalid, maxNoteLength)
	}
	author := rt.Caller().UserID
	if author == "" {
		return nil, fmt.Errorf("%w: a note needs a member behind it", extension.ErrForbidden)
	}
	var added note
	err = rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		subject := "Acme note"
		activity := crm.CreateActivityRequest{
			Kind: crm.CreateActivityRequestKindNote, Source: "manual", Subject: &subject, Body: &body,
		}.LinkTo(crm.CreateActivityRequestLinksEntityTypeDeal, args.DealID)
		if _, err := tx.Core().Activities().Create(ctx, activity); err != nil {
			return err
		}
		added, err = insertNote(ctx, tx, args.DealID, body, "member", &author)
		return err
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(added)
}

// insertNote writes one note with its ledger row and event.
func insertNote(ctx context.Context, tx extension.Tx, dealID, body, kind string, author *string) (note, error) {
	n := note{DealID: dealID, Body: body, Kind: kind, AuthorID: author}
	err := tx.QueryRow(ctx, `
		INSERT INTO ext.ext_acme_greeting_deal_note (deal_id, body, kind, author_id)
		VALUES ($1::uuid, $2, $3, $4::uuid)
		RETURNING id::text, created_at`, dealID, body, kind, author).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return note{}, err
	}
	after, err := json.Marshal(n)
	if err != nil {
		return note{}, err
	}
	payload, err := json.Marshal(map[string]string{"note_id": n.ID, "deal_id": dealID, "kind": kind})
	if err != nil {
		return note{}, err
	}
	return n, tx.Record(ctx,
		extension.Change{Action: extension.AuditCreate, Entity: noteEntity, ID: n.ID, After: after},
		extension.Event{Verb: "note_added", Payload: payload})
}
