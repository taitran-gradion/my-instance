-- SPDX-License-Identifier: BUSL-1.1
-- SPDX-FileCopyrightText: 2026 Gradion
--
-- Acme's notes on a deal. A note from a member names its author; the welcome
-- note the deal decorator files has none, and a deal gets at most one.

CREATE TABLE ext.ext_acme_greeting_deal_note (
    id          uuid        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    deal_id     uuid        NOT NULL,
    body        text        NOT NULL CHECK (length(body) BETWEEN 1 AND 2000),
    kind        text        NOT NULL CHECK (kind IN ('member', 'welcome')),
    author_id   uuid,
    created_at  timestamptz NOT NULL DEFAULT now(),
    CHECK ((kind = 'member') = (author_id IS NOT NULL))
);

CREATE UNIQUE INDEX ext_acme_greeting_deal_note_one_welcome
    ON ext.ext_acme_greeting_deal_note (deal_id) WHERE kind = 'welcome';
CREATE INDEX ext_acme_greeting_deal_note_by_deal
    ON ext.ext_acme_greeting_deal_note (deal_id, created_at DESC);

GRANT SELECT, INSERT ON ext.ext_acme_greeting_deal_note TO margince_app;
