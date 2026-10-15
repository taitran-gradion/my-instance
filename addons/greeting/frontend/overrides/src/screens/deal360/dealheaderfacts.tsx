import { useT } from "@margince/frontend/app";
import { Badge, Row } from "@margince/frontend/design-system";
import { type ComponentProps, useState } from "react";
import { DealSubtitle as CoreDealSubtitle } from "@core-original/src/screens/deal360/dealheaderfacts";
import { AcmeAccountModal } from "../../../../components/account-modal";

// Every other export of core's file passes through unchanged.
export * from "@core-original/src/screens/deal360/dealheaderfacts";

type Props = ComponentProps<typeof CoreDealSubtitle>;

// Core's subtitle, with the Acme mark after it. The deal page hands this the
// whole deal, so the id is read when it is there.
export function DealSubtitle(props: Readonly<{ deal: Props["deal"] & { id?: string } }>) {
  const t = useT();
  const [open, setOpen] = useState(false);
  const dealId = props.deal.id;
  return (
    <Row gap="2" align="center">
      <CoreDealSubtitle deal={props.deal} />
      {dealId && (
        <button type="button" onClick={() => setOpen(true)} style={{ all: "unset", cursor: "pointer" }} aria-haspopup="dialog">
          <Badge tone="info">{t("extAcmeGreeting.dealBadge")}</Badge>
        </button>
      )}
      {open && dealId && <AcmeAccountModal dealId={dealId} onClose={() => setOpen(false)} />}
    </Row>
  );
}
