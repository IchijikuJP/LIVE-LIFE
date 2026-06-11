import { EventCard } from "../components/EventCard";
import type { Copy } from "../i18n";
import type { EventItem, Language } from "../types/api";

// 演出情报：LIVE LIFE 自主演出作为置顶大卡，推荐 / 档案演出在下方网格。
export function ShowsSection({
  copy,
  language,
  ownedEvents,
  recommendedEvents,
}: {
  copy: Copy;
  language: Language;
  ownedEvents: EventItem[];
  recommendedEvents: EventItem[];
}) {
  return (
    <section id="shows" className="section">
      <div className="section-heading split-heading">
        <div>
          <p className="eyebrow">{copy.showsEyebrow}</p>
          <h2>{copy.showsTitle}</h2>
        </div>
        <p>{copy.showsNote}</p>
      </div>

      <div className="event-stack">
        {ownedEvents.map((event) => (
          <EventCard key={event.id} event={event} language={language} copy={copy} featured />
        ))}
      </div>
      <div className="grid">
        {recommendedEvents.map((event) => (
          <EventCard key={event.id} event={event} language={language} copy={copy} />
        ))}
      </div>
    </section>
  );
}
