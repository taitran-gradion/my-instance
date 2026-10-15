import { useT } from "@margince/frontend/app";
import { Badge } from "@margince/frontend/design-system";
import { type ComponentProps, useState } from "react";
import { Fact, RecordFacts } from "@core-original/src/design-system/recordfacts";
import { ContactFacts as CoreContactFacts } from "@core-original/src/screens/contactfacts";
import { AcmeFollowupModal } from "../../../components/followup-modal";

// Every other export of core's file passes through unchanged.
export * from "@core-original/src/screens/contactfacts";

type Contact = ComponentProps<typeof CoreContactFacts>["contact"];

// What Acme needs on a contact before its team reaches out.
function acmeProfile(contact: Contact): { have: number; of: number } {
  const known = [contact.emails?.length, contact.phones?.length, contact.title, contact.employer];
  return { have: known.filter(Boolean).length, of: known.length };
}

// Core's facts strip, then Acme's: how complete the contact is for Acme, and
// its follow-ups.
export function ContactFacts(props: ComponentProps<typeof CoreContactFacts>) {
  const t = useT();
  const [open, setOpen] = useState(false);
  const profile = acmeProfile(props.contact);
  return (
    <div style={{ display: "flex", flexWrap: "wrap", gap: "var(--space-4, 16px)", alignItems: "flex-start" }}>
      <CoreContactFacts {...props} />
      <RecordFacts>
        <Fact label={t("extAcmeGreeting.contact.profile")}>
          <Badge tone={profile.have === profile.of ? "success" : "warning"}>
            {t("extAcmeGreeting.contact.complete", { have: String(profile.have), of: String(profile.of) })}
          </Badge>
        </Fact>
        <Fact label={t("extAcmeGreeting.contact.followups")}>
          <button type="button" onClick={() => setOpen(true)} style={{ all: "unset", cursor: "pointer" }} aria-haspopup="dialog">
            <Badge tone="info">{t("extAcmeGreeting.contact.open")}</Badge>
          </button>
        </Fact>
      </RecordFacts>
      {open && <AcmeFollowupModal contactId={props.contact.id} onClose={() => setOpen(false)} />}
    </div>
  );
}
