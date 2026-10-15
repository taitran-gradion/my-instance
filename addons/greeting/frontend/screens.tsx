import type { CustomScreen } from "@margince/frontend/app";
import { Sparkles } from "lucide-react";
import { AcmeDeskScreen } from "./desk";

// Screens this extension adds to core: each gets #/x/<key>, and `nav` puts it
// in the sidebar group it names.
export const screens: readonly CustomScreen[] = [
  {
    key: "acme-desk",
    component: AcmeDeskScreen,
    nav: {
      group: "records",
      label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" },
      icon: Sparkles,
    },
    palette: { label: { en: "Acme desk", de: "Acme-Desk", vi: "Bàn Acme" } },
  },
];
