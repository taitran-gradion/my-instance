import { problemMessageOf } from "@margince/frontend/api";
import { formatDateTime, useLocale, useT } from "@margince/frontend/app";
import { Badge, Button, Callout, Field, FactList, Row, Select, TextInput } from "@margince/frontend/design-system";
import { useId, useState } from "react";
import { Modal } from "@core-original/src/design-system/modal";
import { refusedForRole, useAcmeFollowups, usePlanAcmeFollowup } from "../api";

const dayChoices = [1, 3, 7, 14] as const;

// Acme's follow-ups on one contact: each is also a task on its timeline.
export function AcmeFollowupModal({ contactId, onClose }: Readonly<{ contactId: string; onClose: () => void }>) {
  const t = useT();
  const { locale } = useLocale();
  const headingId = useId();
  const [note, setNote] = useState("");
  const [days, setDays] = useState("3");
  const followups = useAcmeFollowups(contactId);
  const plan = usePlanAcmeFollowup(contactId);
  const zone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const failure = plan.error ?? followups.error;
  return (
    <Modal open onClose={onClose} labelledBy={headingId} intent="form">
      <h2 id={headingId} className="t-h2 modal-title">
        {t("extAcmeGreeting.followup.title")}
      </h2>
      <div className="form-stack">
        {failure && (
          <Callout tone="warning" title={t("extAcmeGreeting.followup.failed")}>
            {refusedForRole(failure) ? t("extAcmeGreeting.modal.needsGrant") : problemMessageOf(failure, t)}
          </Callout>
        )}
        <FactList
          facts={(followups.data ?? []).map((f) => ({
            key: f.id,
            term: formatDateTime(f.due_at, locale, zone),
            value: (
              <Row gap="2">
                {f.note}
                {f.kind === "welcome" && <Badge tone="info">{t("extAcmeGreeting.followup.automatic")}</Badge>}
              </Row>
            ),
          }))}
        />
        <Field label={t("extAcmeGreeting.followup.note")}>
          {(control) => <TextInput {...control} value={note} onChange={(event) => setNote(event.target.value)} />}
        </Field>
        <Field label={t("extAcmeGreeting.followup.when")}>
          {(control) => (
            <Select
              {...control}
              value={days}
              onChange={setDays}
              options={dayChoices.map((d) => ({ value: String(d), label: t("extAcmeGreeting.followup.inDays", { days: String(d) }) }))}
            />
          )}
        </Field>
        <Row gap="2" justify="end">
          <Button onClick={onClose}>{t("extAcmeGreeting.modal.close")}</Button>
          <Button
            variant="primary"
            disabled={plan.isPending || note.trim() === ""}
            onClick={() => plan.mutate({ days: Number(days), note }, { onSuccess: () => setNote("") })}
          >
            {t("extAcmeGreeting.followup.plan")}
          </Button>
        </Row>
      </div>
    </Modal>
  );
}
