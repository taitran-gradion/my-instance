import { api } from "@margince/frontend/api";
import { useT } from "@margince/frontend/app";
import { Badge, Callout, Card, SectionHeader } from "@margince/frontend/design-system";
import { useQuery } from "@tanstack/react-query";

// #/ext/acme-greeting: this installation's own screen, read from the instance
// repository rather than from a published module.
export function AcmeGreetingPage() {
  const t = useT();
  const greeting = useQuery({
    queryKey: ["ext", "acme-greeting", "hello"],
    queryFn: async () => {
      const { data, error } = await api.GET("/ext/acme-greeting/hello");
      if (error) throw error;
      return data;
    },
  });
  return (
    <div className="wrap narrow">
      <SectionHeader title={t("extAcmeGreeting.title")} level={1} />
      <Callout tone="accent" title={t("extAcmeGreeting.local.title")}>
        {t("extAcmeGreeting.local.body")}
      </Callout>
      <Card>
        <SectionHeader title={t("extAcmeGreeting.route.title")} level={2} />
        <p>
          {greeting.data?.message ?? (
            <Badge tone={greeting.isError ? "danger" : "info"}>
              {greeting.isError ? t("extAcmeGreeting.route.failed") : t("extAcmeGreeting.route.asking")}
            </Badge>
          )}
        </p>
      </Card>
    </div>
  );
}
