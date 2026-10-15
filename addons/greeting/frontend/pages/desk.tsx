import { useT } from "@margince/frontend/app";
import { Badge, Callout, Card, DataTable, EmptyState, SectionHeader } from "@margince/frontend/design-system";
import { useAcmeNotes } from "../api";

type NoteRow = Readonly<{ id: string; deal_id: string; body: string; kind: string }>;

// #/x/acme-desk: the account team's latest notes across deals.
export function AcmeDeskPage() {
  const t = useT();
  const notes = useAcmeNotes();
  const rows: NoteRow[] = notes.data ?? [];
  return (
    <div className="wrap">
      <Callout tone="info" title={t("extAcmeGreeting.desk.how.title")}>
        {t("extAcmeGreeting.desk.how.body")}
      </Callout>
      <Card>
        <SectionHeader title={t("extAcmeGreeting.desk.notes")} level={2} />
        {rows.length === 0 ? (
          <EmptyState>{t("extAcmeGreeting.desk.empty")}</EmptyState>
        ) : (
          <DataTable<NoteRow>
            label={t("extAcmeGreeting.desk.notes")}
            rows={rows}
            rowKey={(row) => row.id}
            onRowClick={(row) => {
              window.location.hash = `#/deals/${row.deal_id}`;
            }}
            columns={[
              { key: "body", header: t("extAcmeGreeting.desk.col.note"), render: (row) => row.body },
              {
                key: "kind",
                header: t("extAcmeGreeting.desk.col.kind"),
                render: (row) => (
                  <Badge tone="info">
                    {row.kind === "welcome" ? t("extAcmeGreeting.note.welcome") : t("extAcmeGreeting.note.member")}
                  </Badge>
                ),
              },
            ]}
          />
        )}
      </Card>
    </div>
  );
}
