import { api } from "@margince/frontend/api";
import { useT } from "@margince/frontend/app";
import { Callout, Card, DataTable, EmptyState, SectionHeader } from "@margince/frontend/design-system";
import { useQuery } from "@tanstack/react-query";
import { Sparkles } from "lucide-react";
import type { CustomScreen } from "@core-original/src/app/custom";
import { EntityRef } from "@core-original/src/components/entityref";

type DealRow = Readonly<{ id: string; name: string; company_id?: string | null }>;

// #/x/acme-desk: a screen core does not have, added to core's custom/ folder by
// this installation's own extension, so the sidebar lists it under Records.
function AcmeDeskScreen() {
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
                key: "company",
                header: t("extAcmeGreeting.desk.col.company"),
                render: (row) => (row.company_id ? <EntityRef kind="company" id={row.company_id} /> : "—"),
              },
            ]}
          />
        )}
      </Card>
    </div>
  );
}

export const screen: CustomScreen = {
  key: "acme-desk",
  component: AcmeDeskScreen,
  nav: {
    group: "records",
    label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" },
    icon: Sparkles,
  },
  palette: { label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" } },
};
