// SPDX-License-Identifier: BUSL-1.1
// SPDX-FileCopyrightText: 2026 Gradion

package acmegreeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/margince/margince/backend/pkg/extension"
	"github.com/margince/margince/backend/pkg/extension/crm"
)

const (
	followupEntity     = "ext_acme_greeting_contact_followup"
	maxFollowupNote    = 500
	maxFollowupDays    = 30
	introduceAfterDays = 2
)

type followup struct {
	ID         string    `json:"id"`
	ContactID  string    `json:"contact_id"`
	ActivityID string    `json:"activity_id"`
	Note       string    `json:"note"`
	Kind       string    `json:"kind"`
	DueAt      time.Time `json:"due_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// titleCaseNames runs before core's createContact and updateContact: Acme
// keeps names capitalized, so "trab tran" is stored as "Trab Tran".
func titleCaseNames(_ context.Context, _ extension.Runtime, op *extension.Operation) error {
	var body map[string]any
	if err := json.Unmarshal(op.Body, &body); err != nil {
		return fmt.Errorf("%w: the contact is not a JSON object", extension.ErrInvalid)
	}
	changed := false
	for _, field := range []string{"full_name", "first_name", "last_name"} {
		if name, ok := body[field].(string); ok && titleCase(name) != name {
			body[field] = titleCase(name)
			changed = true
		}
	}
	if !changed {
		return nil
	}
	out, err := json.Marshal(body)
	if err != nil {
		return err
	}
	op.Body = out
	return nil
}

// titleCase capitalizes the first letter of each word and leaves the rest as
// typed, so "mcDonald" keeps its inner capital.
func titleCase(name string) string {
	words := strings.Fields(name)
	for i, w := range words {
		r, size := utf8.DecodeRuneInString(w)
		words[i] = string(unicode.ToUpper(r)) + w[size:]
	}
	return strings.Join(words, " ")
}

// introduceAcme runs inside the transaction that creates a contact: it plans
// an "Introduce Acme" task, so a new contact never lands without a next step.
// A contact created with nobody behind it (an import job) gets none, since a
// core task needs a member to be filed for.
func introduceAcme(ctx context.Context, rt extension.Runtime, w extension.Write) error {
	if rt.Caller().UserID == "" {
		return nil
	}
	return rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		_, err := planFollowupIn(ctx, tx, rt.Caller().UserID, w.ID, "Introduce Acme", introduceAfterDays, "welcome")
		return err
	})
}

// listFollowups answers GET /ext/acme-greeting/followups for one contact.
func listFollowups(ctx context.Context, rt extension.Runtime, in json.RawMessage) (json.RawMessage, error) {
	args, err := extension.DecodeArgs[struct {
		ContactID string `json:"contact_id"`
	}](in)
	if err != nil {
		return nil, err
	}
	if !extension.IsCanonicalUUID(args.ContactID) {
		return nil, fmt.Errorf("%w: contact_id is not a contact id", extension.ErrInvalid)
	}
	out := []followup{}
	err = rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id::text, contact_id::text, activity_id::text, note, kind, due_at, created_at
			FROM ext.ext_acme_greeting_contact_followup
			WHERE contact_id = $1::uuid
			ORDER BY due_at`, args.ContactID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var f followup
			if err := rows.Scan(&f.ID, &f.ContactID, &f.ActivityID, &f.Note, &f.Kind, &f.DueAt, &f.CreatedAt); err != nil {
				return err
			}
			out = append(out, f)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string][]followup{"followups": out})
}

// planFollowup answers POST /ext/acme-greeting/followups. The core task is
// filed through the timeline door, which checks the caller may see the contact.
func planFollowup(ctx context.Context, rt extension.Runtime, in json.RawMessage) (json.RawMessage, error) {
	args, err := extension.DecodeArgs[struct {
		ContactID string `json:"contact_id"`
		Days      int    `json:"days"`
		Note      string `json:"note"`
	}](in)
	if err != nil {
		return nil, err
	}
	note := strings.TrimSpace(args.Note)
	switch {
	case !extension.IsCanonicalUUID(args.ContactID):
		return nil, fmt.Errorf("%w: contact_id is not a contact id", extension.ErrInvalid)
	case args.Days < 1 || args.Days > maxFollowupDays:
		return nil, fmt.Errorf("%w: a follow-up is 1 to %d days out", extension.ErrInvalid, maxFollowupDays)
	case note == "" || len(note) > maxFollowupNote:
		return nil, fmt.Errorf("%w: a follow-up note is 1 to %d characters", extension.ErrInvalid, maxFollowupNote)
	}
	member := rt.Caller().UserID
	if member == "" {
		return nil, fmt.Errorf("%w: a follow-up needs a member to be filed for", extension.ErrForbidden)
	}
	var planned followup
	err = rt.Tx(ctx, func(ctx context.Context, tx extension.Tx) error {
		planned, err = planFollowupIn(ctx, tx, member, args.ContactID, note, args.Days, "member")
		return err
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(planned)
}

// planFollowupIn files the core task and Acme's row for it, with the row's
// ledger entry, on one transaction.
func planFollowupIn(ctx context.Context, tx extension.Tx, member, contactID, note string, days int, kind string) (followup, error) {
	due := time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour).Truncate(time.Hour)
	subject := "Acme follow-up: " + note
	task, err := tx.Core().Activities().Create(ctx, crm.CreateActivityRequest{
		Kind: crm.CreateActivityRequestKindTask, Source: "manual", Subject: &subject, DueAt: &due, AssigneeId: &member,
	}.LinkTo(crm.CreateActivityRequestLinksEntityTypeContact, contactID))
	if err != nil {
		return followup{}, err
	}
	f := followup{ContactID: contactID, ActivityID: task.Id, Note: note, Kind: kind, DueAt: due}
	err = tx.QueryRow(ctx, `
		INSERT INTO ext.ext_acme_greeting_contact_followup (contact_id, activity_id, note, kind, due_at, created_by)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6::uuid)
		RETURNING id::text, created_at`, contactID, task.Id, note, kind, due, member).Scan(&f.ID, &f.CreatedAt)
	if err != nil {
		return followup{}, err
	}
	after, err := json.Marshal(f)
	if err != nil {
		return followup{}, err
	}
	payload, err := json.Marshal(map[string]string{"followup_id": f.ID, "contact_id": contactID, "kind": kind})
	if err != nil {
		return followup{}, err
	}
	return f, tx.Record(ctx,
		extension.Change{Action: extension.AuditCreate, Entity: followupEntity, ID: f.ID, After: after},
		extension.Event{Verb: "followup_planned", Payload: payload})
}
