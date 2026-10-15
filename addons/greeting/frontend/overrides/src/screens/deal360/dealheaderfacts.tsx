import { api } from "@margince/frontend/api";
import { useT } from "@margince/frontend/app";
import { Badge, Button, Callout, Field, FactList, Row, TextInput } from "@margince/frontend/design-system";
import { useMutation } from "@tanstack/react-query";
import { type ComponentProps, useId, useState } from "react";
import { EntityRef } from "@core-original/src/components/entityref";
import { Modal } from "@core-original/src/design-system/modal";
import { DealSubtitle as CoreDealSubtitle } from "@core-original/src/screens/deal360/dealheaderfacts";

// Every other export of core's file passes through unchanged.
export * from "@core-original/src/screens/deal360/dealheaderfacts";

// Core's subtitle, with this installation's mark after it; the mark opens the
// account's Acme view.
export function DealSubtitle(props: ComponentProps<typeof CoreDealSubtitle>) {
  const t = useT();
  const [open, setOpen] = useState(false);
  return (
    <Row gap="2" align="center">
      <CoreDealSubtitle {...props} />
      <button
        type="button"
        onClick={() => setOpen(true)}
        style={{ all: "unset", cursor: "pointer" }}
        aria-haspopup="dialog"
      >
        <Badge tone="info">{t("extAcmeGreeting.dealBadge")}</Badge>
      </button>
      {open && <AcmeAccountModal companyId={props.deal.company_id ?? null} onClose={() => setOpen(false)} />}
    </Row>
  );
}

function AcmeAccountModal({ companyId, onClose }: Readonly<{ companyId: string | null; onClose: () => void }>) {
  const t = useT();
  const headingId = useId();
  const [note, setNote] = useState("");
  const ask = useMutation({
    mutationFn: async () => {
      const { data, error } = await api.GET("/ext/acme-greeting/hello");
      if (error) throw error;
      return data;
    },
  });
  return (
    <Modal open onClose={onClose} labelledBy={headingId} intent="form">
      <h2 id={headingId} className="t-h2 modal-title">
        {t("extAcmeGreeting.modal.title")}
      </h2>
      <div className="form-stack">
        <Callout tone="info" title={t("extAcmeGreeting.modal.how.title")}>
          {t("extAcmeGreeting.modal.how.body")}
        </Callout>
        <FactList
          facts={[
            {
              key: "company",
              term: t("extAcmeGreeting.modal.company"),
              value: companyId ? <EntityRef kind="company" id={companyId} /> : "—",
            },
            {
              key: "answer",
              term: t("extAcmeGreeting.modal.answer"),
              value: ask.data?.message ?? (ask.isError ? t("extAcmeGreeting.route.failed") : "—"),
            },
          ]}
        />
        <Field label={t("extAcmeGreeting.modal.note")}>
          {(control) => <TextInput {...control} value={note} onChange={(event) => setNote(event.target.value)} />}
        </Field>
        <Row gap="2" justify="end">
          <Button onClick={onClose}>{t("extAcmeGreeting.modal.close")}</Button>
          <Button variant="primary" disabled={ask.isPending} onClick={() => ask.mutate()}>
            {ask.isPending ? t("extAcmeGreeting.route.asking") : t("extAcmeGreeting.modal.ask")}
          </Button>
        </Row>
      </div>
    </Modal>
  );
}
