-- SPDX-License-Identifier: BUSL-1.1
-- SPDX-FileCopyrightText: 2026 Gradion
--
-- Acme's follow-ups on a contact. Each is also a core task on the contact's
-- timeline (activity_id); this row keeps why Acme planned it.

CREATE TABLE ext.ext_acme_greeting_contact_followup (
    id           uuid        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    contact_id   uuid        NOT NULL,
    activity_id  uuid        NOT NULL,
    note         text        NOT NULL CHECK (length(note) BETWEEN 1 AND 500),
    kind         text        NOT NULL CHECK (kind IN ('member', 'welcome')),
    due_at       timestamptz NOT NULL,
    created_by   uuid        NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ext_acme_greeting_contact_followup_by_contact
    ON ext.ext_acme_greeting_contact_followup (contact_id, due_at);

GRANT SELECT, INSERT ON ext.ext_acme_greeting_contact_followup TO margince_app;
