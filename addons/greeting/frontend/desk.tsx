import { api } from "@margince/frontend/api";
import { useT } from "@margince/frontend/app";
import { Badge, Callout, Card, DataTable, EmptyState, SectionHeader } from "@margince/frontend/design-system";
import { useQuery } from "@tanstack/react-query";

type DealRow = Readonly<{ id: string; name: string; status: string }>;

// #/x/acme-desk: this installation's own page, listed in the sidebar after Deals.
export function AcmeDeskScreen() {
  const t = useT();
  const deals = useQuery({
    queryKey: ["ext", "acme-greeting", "desk", "deals"],
    queryFn: async () => {
      const { data, error } = await api.GET("/deals", { params: { query: { limit: 20 } } });
      if (error) throw error;
      return data.data;
    },
  });
  const rows: DealRow[] = deals.data ?? [];
  return (
    <div className="wrap">
      <Callout tone="info" title={t("extAcmeGreeting.desk.how.title")}>
        {t("extAcmeGreeting.desk.how.body")}
      </Callout>
      <Card>
        <SectionHeader title={t("extAcmeGreeting.desk.deals")} level={2} />
        {rows.length === 0 ? (
          <EmptyState>{t("extAcmeGreeting.desk.empty")}</EmptyState>
        ) : (
          <DataTable<DealRow>
            label={t("extAcmeGreeting.desk.deals")}
            rows={rows}
            rowKey={(row) => row.id}
            onRowClick={(row) => {
              window.location.hash = `#/deals/${row.id}`;
            }}
            columns={[
              { key: "name", header: t("extAcmeGreeting.desk.col.name"), render: (row) => row.name },
              {
                key: "status",
                header: t("extAcmeGreeting.desk.col.status"),
                render: (row) => <Badge tone="info">{row.status}</Badge>,
              },
            ]}
          />
        )}
      </Card>
    </div>
  );
}
