import type { Copy } from "../i18n";
import { localized } from "../lib/localized";
import type { CatalogItem, EventItem, Language } from "../types/api";

// 首页 Schedule 面板里的一行，演出和 CD 单品共用。靠是否有 date 字段区分两类。
export function ScheduleRow({ item, language, copy }: { item: EventItem | CatalogItem; language: Language; copy: Copy }) {
  const isEvent = "date" in item;
  return (
    <article className="schedule-item">
      <div className="schedule-date">{isEvent ? item.date : "TBD"}</div>
      <div>
        <span>{isEvent ? copy.scheduleLive : item.format === "vinyl" ? copy.scheduleVinyl : copy.scheduleCD}</span>
        <strong>{localized(item, "title", language)}</strong>
        <small>{isEvent ? item.venue : `${item.artist} / ${item.status || item.price}`}</small>
      </div>
    </article>
  );
}
