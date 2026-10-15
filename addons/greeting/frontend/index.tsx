import type { CustomScreen } from "@margince/frontend/app";
import { Sparkles } from "lucide-react";
import { AcmeDeskPage } from "./pages/desk";
import { AcmeGreetingPage } from "./pages/greeting";

// Everything this extension's frontend gives core is named here.

// The extension's own page, at #/ext/acme-greeting.
export default AcmeGreetingPage;

// Pages added beside core's: each gets #/x/<key>, and `nav` a sidebar row.
export const screens: readonly CustomScreen[] = [
  {
    key: "acme-desk",
    component: AcmeDeskPage,
    nav: {
      group: "records",
      label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" },
      icon: Sparkles,
    },
    palette: { label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" } },
  },
];
