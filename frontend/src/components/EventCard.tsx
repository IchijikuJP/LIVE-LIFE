import type { Copy } from "../i18n";
import { localized } from "../lib/localized";
import type { EventItem, Language } from "../types/api";

// 演出卡片。featured=true 用于 LIVE LIFE 自主演出（置顶大卡），否则是推荐 / 档案演出。
export function EventCard({ event, language, copy, featured = false }: { event: EventItem; language: Language; copy: Copy; featured?: boolean }) {
  return (
    <article className={`event-card ${featured ? "featured" : ""}`}>
      {event.imageUrl ? (
        <img className="event-image" src={event.imageUrl} alt={localized(event, "title", language)} />
      ) : (
        <div className="event-image placeholder" aria-hidden="true" />
      )}
      <div className="event-body">
        <span className="pill">{featured ? copy.ownBadge : copy.recommendBadge}</span>
        <h3>{localized(event, "title", language)}</h3>
        <div className="meta">
          <span>{event.date}</span>
          <span>{event.time}</span>
          <span>{event.venue}</span>
          <span>{event.price}</span>
        </div>
        <p className="summary">{localized(event, "summary", language)}</p>
        <div className="lineup">{event.lineup?.map((name) => <span key={name}>{name}</span>)}</div>
        <p className="note"><strong>{copy.ticketPending}</strong> {localized(event, "ticketNote", language) || copy.ticketPending}</p>
        {localized(event, "sourceNote", language) ? (
          <p className="source-note"><strong>{copy.sourceNote}</strong> {localized(event, "sourceNote", language)}</p>
        ) : null}
        <div className="tags">{event.tags?.map((tag) => <span className="tag" key={tag}>{tag}</span>)}</div>
      </div>
    </article>
  );
}
