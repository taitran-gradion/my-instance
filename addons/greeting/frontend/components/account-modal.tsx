import { problemMessageOf } from "@margince/frontend/api";
import { useT } from "@margince/frontend/app";
import { Badge, Button, Callout, Field, FactList, Row, TextInput } from "@margince/frontend/design-system";
import { useId, useState } from "react";
import { Modal } from "@core-original/src/design-system/modal";
import { refusedForRole, useAcmeNotes, useAddAcmeNote, useDealScore } from "../api";

// The Acme view of one deal: the score this extension provides through core,
// and the notes it keeps in its own table.
export function AcmeAccountModal({ dealId, onClose }: Readonly<{ dealId: string; onClose: () => void }>) {
  const t = useT();
  const headingId = useId();
  const [draft, setDraft] = useState("");
  const score = useDealScore(dealId);
  const notes = useAcmeNotes(dealId);
  const add = useAddAcmeNote(dealId);
  return (
    <Modal open onClose={onClose} labelledBy={headingId} intent="form">
      <h2 id={headingId} className="t-h2 modal-title">
        {t("extAcmeGreeting.modal.title")}
      </h2>
      <div className="form-stack">
        <FactList
          facts={[
            {
              key: "score",
              term: t("extAcmeGreeting.modal.score"),
              value: score.data ? <Badge tone="info">{score.data.score} / 100</Badge> : "—",
              note: score.data?.reason,
            },
          ]}
        />
        {(notes.isError || add.isError) && (
          <Callout tone="warning" title={t(add.isError ? "extAcmeGreeting.modal.addFailed" : "extAcmeGreeting.modal.notesFailed")}>
            {refusedForRole(add.error ?? notes.error)
              ? t("extAcmeGreeting.modal.needsGrant")
              : problemMessageOf(add.error ?? notes.error, t)}
          </Callout>
        )}
        <FactList
          facts={(notes.data ?? []).map((note) => ({
            key: note.id,
            term: note.kind === "welcome" ? t("extAcmeGreeting.note.welcome") : t("extAcmeGreeting.note.member"),
            value: note.body,
          }))}
        />
        <Field label={t("extAcmeGreeting.modal.note")}>
          {(control) => <TextInput {...control} value={draft} onChange={(event) => setDraft(event.target.value)} />}
        </Field>
        <Row gap="2" justify="end">
          <Button onClick={onClose}>{t("extAcmeGreeting.modal.close")}</Button>
          <Button
            variant="primary"
            disabled={add.isPending || draft.trim() === ""}
            onClick={() => add.mutate(draft, { onSuccess: () => setDraft("") })}
          >
            {t("extAcmeGreeting.modal.add")}
          </Button>
        </Row>
      </div>
    </Modal>
  );
}
