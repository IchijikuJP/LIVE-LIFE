import { ScheduleRow } from "../components/ScheduleRow";
import { V2Texture } from "../design/V2Texture";
import { V3Signal } from "../design/V3Signal";
import type { Copy } from "../i18n";
import type { ApiStatus } from "../hooks/useLiveLifeData";
import type { CatalogItem, EventItem, HealthPayload, Language } from "../types/api";

function apiStatusLabel(status: ApiStatus, copy: Copy) {
  if (status === "online") return copy.apiOnline;
  if (status === "offline") return copy.apiOffline;
  return copy.apiChecking;
}

// 首页主视觉：设计纹理 + 品牌标题 + 右侧 Schedule 面板（含 API 状态、近期演出与 CD）。
export function HeroSection({
  copy,
  language,
  apiStatus,
  health,
  ownedEvents,
  catalog,
}: {
  copy: Copy;
  language: Language;
  apiStatus: ApiStatus;
  health: HealthPayload | null;
  ownedEvents: EventItem[];
  catalog: CatalogItem[];
}) {
  return (
    <section id="home" className="hero-grid">
      <div className="hero-stage">
        <V2Texture />
        <V3Signal />
        <div className="hero-lockup">
          <p className="eyebrow">{copy.heroEyebrow}</p>
          <h1>{copy.heroTitle}</h1>
          <p className="lead">{copy.heroLead}</p>
          <div className="hero-actions">
            <a className="button primary" href="#shows">{copy.heroPrimary}</a>
            <a className="button secondary" href="#cd-select">{copy.heroSecondary}</a>
          </div>
        </div>
      </div>

      <aside className="schedule-panel" aria-labelledby="scheduleTitle">
        <div className="panel-topline">
          <span className={`status-dot ${apiStatus === "online" ? "ok" : ""}`} />
          <div>
            <strong>{apiStatusLabel(apiStatus, copy)}</strong>
            <span>{health ? `${health.brand} / ${health.service}` : "/api/health"}</span>
          </div>
        </div>
        <p className="eyebrow">{copy.scheduleEyebrow}</p>
        <h2 id="scheduleTitle">{copy.scheduleTitle}</h2>
        <div className="schedule-list">
          {[...ownedEvents.slice(0, 2), ...catalog.slice(0, 2)].map((item) => (
            <ScheduleRow key={item.id} item={item} language={language} copy={copy} />
          ))}
        </div>
      </aside>
    </section>
  );
}
