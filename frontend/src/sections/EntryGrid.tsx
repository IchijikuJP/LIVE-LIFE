import { Archive, Disc3, Mail, Ticket } from "lucide-react";

import { Entry } from "../components/Entry";
import type { Copy } from "../i18n";

// 四个板块入口磁贴网格。
export function EntryGrid({ copy }: { copy: Copy }) {
  return (
    <section className="entry-grid" aria-label="LIVE LIFE sections">
      <Entry href="#shows" index="01" label={copy.navShows} icon={<Ticket size={20} />} />
      <Entry href="#cd-select" index="02" label={copy.navCDSelect} icon={<Disc3 size={20} />} />
      <Entry href="#archive" index="03" label={copy.navArchive} icon={<Archive size={20} />} />
      <Entry href="#connect" index="04" label={copy.navConnect} icon={<Mail size={20} />} />
    </section>
  );
}
