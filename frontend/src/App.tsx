import { Archive, Disc3, Mail, Ticket } from "lucide-react";
import { FormEvent, ReactNode, useEffect, useMemo, useState } from "react";

import type { CatalogItem, EventItem, Language, LocalizedText } from "./types/api";
import { copy, languageLabels } from "./i18n";
import type { Copy } from "./i18n";
import { useLiveLifeData } from "./hooks/useLiveLifeData";

// UI 本地状态类型，不属于 API 契约，后续随设计/筛选功能再各自归位。
type DesignVariant = "v2" | "v2-refined" | "v3";
type FormatFilter = "all" | "cd" | "vinyl";

const connectMailAddress = "livelife.cn.2023@gmail.com";

function formText(payload: Record<string, FormDataEntryValue>, key: string) {
  const value = payload[key];
  return typeof value === "string" ? value.trim() : "";
}

function connectTopicLabel(topic: string, copyText: Copy) {
  const labels: Record<string, string> = {
    ticket: copyText.topicTicket,
    "cd-select": copyText.topicCDSelect,
    support: copyText.topicSupport,
    collab: copyText.topicCollab,
  };
  return labels[topic] || topic;
}

// Connect 表单仍然先使用浏览器原生 mailto 能力，用户点击后会打开自己的邮件 App 或网页邮箱草稿。
// 同一份内容还会继续 POST 到本地 API，方便之后接后台列表、客服系统或真正的服务端邮件发送。
function buildConnectMailto(payload: Record<string, FormDataEntryValue>, copyText: Copy) {
  const nickname = formText(payload, "nickname") || "-";
  const email = formText(payload, "email") || "-";
  const topic = connectTopicLabel(formText(payload, "topic"), copyText);
  const message = formText(payload, "message") || "-";
  const subject = `[LIVE LIFE CONNECT] ${topic} / ${nickname}`;
  const body = [
    "LIVE LIFE CONNECT",
    "",
    `${copyText.labelNickname}: ${nickname}`,
    `${copyText.labelEmail}: ${email}`,
    `${copyText.labelTopic}: ${topic}`,
    "",
    `${copyText.labelMessage}:`,
    message,
    "",
    `Page: ${window.location.href}`,
    `Time: ${new Date().toLocaleString()}`,
  ].join("\n");

  return `mailto:${connectMailAddress}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`;
}

function openConnectMail(payload: Record<string, FormDataEntryValue>, copyText: Copy) {
  const mailtoUrl = buildConnectMailto(payload, copyText);
  const link = document.createElement("a");
  link.href = mailtoUrl;
  link.target = "_blank";
  link.rel = "noopener noreferrer";
  document.body.appendChild(link);
  link.click();
  link.remove();
}


const classicBandCloud =
  "THE BEATLES THE ROLLING STONES LED ZEPPELIN PINK FLOYD QUEEN THE WHO DAVID BOWIE THE KINKS T. REX ROXY MUSIC THE CLASH SEX PISTOLS THE JAM BUZZCOCKS JOY DIVISION NEW ORDER THE SMITHS THE CURE SIOUXSIE AND THE BANSHEES ECHO AND THE BUNNYMEN THE STONE ROSES MY BLOODY VALENTINE SLOWDIVE PRIMAL SCREAM RADIOHEAD OASIS BLUR SUEDE PULP THE VERVE ARCTIC MONKEYS NIRVANA PIXIES SONIC YOUTH R.E.M. THE STOOGES THE VELVET UNDERGROUND TALKING HEADS";

const v2BandNames = [
  "THE BEATLES",
  "THE ROLLING STONES",
  "THE KINKS",
  "THE WHO",
  "THE YARDBIRDS",
  "CREAM",
  "LED ZEPPELIN",
  "PINK FLOYD",
  "BLACK SABBATH",
  "DEEP PURPLE",
  "QUEEN",
  "DAVID BOWIE",
  "T. REX",
  "ROXY MUSIC",
  "SPARKS",
  "THE VELVET UNDERGROUND",
  "THE STOOGES",
  "MC5",
  "NEW YORK DOLLS",
  "PATTI SMITH",
  "TELEVISION",
  "RAMONES",
  "BLONDIE",
  "TALKING HEADS",
  "THE CLASH",
  "SEX PISTOLS",
  "BUZZCOCKS",
  "WIRE",
  "XTC",
  "THE JAM",
  "PUBLIC IMAGE LTD",
  "JOY DIVISION",
  "NEW ORDER",
  "THE FALL",
  "GANG OF FOUR",
  "SIOUXSIE AND THE BANSHEES",
  "BAUHAUS",
  "THE CURE",
  "ECHO AND THE BUNNYMEN",
  "THE PSYCHEDELIC FURS",
  "THE SMITHS",
  "THE JESUS AND MARY CHAIN",
  "MY BLOODY VALENTINE",
  "SLOWDIVE",
  "RIDE",
  "LUSH",
  "COCTEAU TWINS",
  "SPACEMEN 3",
  "THE STONE ROSES",
  "HAPPY MONDAYS",
  "PRIMAL SCREAM",
  "MASSIVE ATTACK",
  "PORTISHEAD",
  "TRICKY",
  "RADIOHEAD",
  "OASIS",
  "BLUR",
  "SUEDE",
  "PULP",
  "THE VERVE",
  "MOGWAI",
  "PJ HARVEY",
  "NIRVANA",
  "PEARL JAM",
  "SONIC YOUTH",
  "PIXIES",
  "DINOSAUR JR.",
  "R.E.M.",
  "THE REPLACEMENTS",
  "HUSKER DU",
  "GUIDED BY VOICES",
  "PAVEMENT",
  "YO LA TENGO",
  "THE FLAMING LIPS",
  "WILCO",
  "THE WHITE STRIPES",
  "THE STROKES",
  "INTERPOL",
  "YEAH YEAH YEAHS",
  "LCD SOUNDSYSTEM",
  "ARCADE FIRE",
  "TV ON THE RADIO",
  "ARCTIC MONKEYS",
  "FRANZ FERDINAND",
  "BLOC PARTY",
  "THE LIBERTINES",
  "FOALS",
  "IDLES",
  "FONTAINES D.C.",
  "BLACK MIDI",
  "SQUID",
  "CAN",
  "NEU!",
  "KRAFTWERK",
  "FAUST",
  "TALK TALK",
  "WIRE",
  "SWANS",
  "THE RESIDENTS",
  "SUICIDE",
  "DEVO",
  "MISSION OF BURMA",
  "MINUTEMEN",
  "FUGAZI",
  "BIG BLACK",
  "SLINT",
  "TORTOISE",
  "STEREOLAB",
  "BROADCAST",
  "THE SEA AND CAKE",
];


function readInitialDesign(): DesignVariant {
  const params = new URLSearchParams(window.location.search);
  const fromUrl = params.get("design");
  if (fromUrl === "v2" || fromUrl === "v2-refined" || fromUrl === "v3") {
    return fromUrl;
  }
  const stored = window.localStorage.getItem("liveLifeDesignVariant");
  return stored === "v3" || stored === "v2-refined" ? stored : "v2";
}

export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [design, setDesign] = useState<DesignVariant>(readInitialDesign);
  const [formatFilter, setFormatFilter] = useState<FormatFilter>("all");
  const [formStatus, setFormStatus] = useState("");
  const { events, catalog, contents, health, apiStatus } = useLiveLifeData();
  const t = copy[language];

  const ownedEvents = useMemo(() => events.filter((event) => event.owned), [events]);
  const recommendedEvents = useMemo(() => events.filter((event) => !event.owned), [events]);
  const filteredCatalog = useMemo(
    () => (formatFilter === "all" ? catalog : catalog.filter((item) => item.format === formatFilter)),
    [catalog, formatFilter],
  );

  // 页面语言不只是按钮状态，也会写入 html lang，方便浏览器、翻译工具和无障碍工具正确识别。
  useEffect(() => {
    document.documentElement.lang = t.langAttr;
  }, [t.langAttr]);

  // 设计方案只改变前端视觉，不改变 API 路由和数据字段。URL 参数保留 design，方便把具体方案链接发给客户 review。
  useEffect(() => {
    document.body.dataset.design = design;
    window.localStorage.setItem("liveLifeDesignVariant", design);
    const url = new URL(window.location.href);
    url.searchParams.set("design", design);
    url.searchParams.set("v", "20260611-v2v3-logo-mail");
    window.history.replaceState({}, "", url);
  }, [design]);

  async function submitConnect(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormStatus(t.submitting);

    const form = event.currentTarget;
    const payload = Object.fromEntries(new FormData(form).entries()) as Record<string, FormDataEntryValue>;
    openConnectMail(payload, t);

    try {
      const response = await fetch("/api/connect", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const result = await response.json();
      if (!response.ok) {
        throw new Error(result.error || t.submitFailed);
      }
      setFormStatus(result.message || t.submitFallback);
      form.reset();
    } catch (error) {
      setFormStatus(error instanceof Error ? error.message : t.submitFailed);
    }
  }

  return (
    <div className="app-shell" data-design={design}>
      <header className="topbar">
        <a className="brand" href="#home" aria-label="LIVE LIFE">
          <span className="brand-mark" aria-hidden="true" />
          <span>LIVE LIFE</span>
        </a>

        <div className="topbar-actions">
          <nav aria-label="Primary navigation">
            <a href="#shows">{t.navShows}</a>
            <a href="#cd-select">{t.navCDSelect}</a>
            <a href="#archive">{t.navArchive}</a>
            <a href="#connect">{t.navConnect}</a>
          </nav>

          <label className="design-switcher">
            <span>{t.designLabel}</span>
            <select value={design} onChange={(event) => setDesign(event.target.value as DesignVariant)}>
              <option value="v2">{t.designV2}</option>
              <option value="v2-refined">{t.designV2Refined}</option>
              <option value="v3">{t.designV3}</option>
            </select>
          </label>

          <div className="language-switcher" aria-label="Language">
            {(["zh", "ja", "en"] as const).map((item) => (
              <button
                className={language === item ? "active" : ""}
                key={item}
                type="button"
                aria-pressed={language === item}
                onClick={() => setLanguage(item)}
              >
                {languageLabels[item]}
              </button>
            ))}
          </div>
        </div>
      </header>

      <main>
        <section id="home" className="hero-grid">
          <div className="hero-stage">
            <V2Texture />
            <V3Signal />
            <div className="hero-lockup">
              <p className="eyebrow">{t.heroEyebrow}</p>
              <h1>{t.heroTitle}</h1>
              <p className="lead">{t.heroLead}</p>
              <div className="hero-actions">
                <a className="button primary" href="#shows">{t.heroPrimary}</a>
                <a className="button secondary" href="#cd-select">{t.heroSecondary}</a>
              </div>
            </div>
          </div>

          <aside className="schedule-panel" aria-labelledby="scheduleTitle">
            <div className="panel-topline">
              <span className={`status-dot ${apiStatus === "online" ? "ok" : ""}`} />
              <div>
                <strong>{apiStatusLabel(apiStatus, t)}</strong>
                <span>{health ? `${health.brand} / ${health.service}` : "/api/health"}</span>
              </div>
            </div>
            <p className="eyebrow">{t.scheduleEyebrow}</p>
            <h2 id="scheduleTitle">{t.scheduleTitle}</h2>
            <div className="schedule-list">
              {[...ownedEvents.slice(0, 2), ...catalog.slice(0, 2)].map((item) => (
                <ScheduleRow key={item.id} item={item} language={language} copy={t} />
              ))}
            </div>
          </aside>
        </section>

        <section className="entry-grid" aria-label="LIVE LIFE sections">
          <Entry href="#shows" index="01" label={t.navShows} icon={<Ticket size={20} />} />
          <Entry href="#cd-select" index="02" label={t.navCDSelect} icon={<Disc3 size={20} />} />
          <Entry href="#archive" index="03" label={t.navArchive} icon={<Archive size={20} />} />
          <Entry href="#connect" index="04" label={t.navConnect} icon={<Mail size={20} />} />
        </section>

        <section id="shows" className="section">
          <div className="section-heading split-heading">
            <div>
              <p className="eyebrow">{t.showsEyebrow}</p>
              <h2>{t.showsTitle}</h2>
            </div>
            <p>{t.showsNote}</p>
          </div>

          <div className="event-stack">
            {ownedEvents.map((event) => (
              <EventCard key={event.id} event={event} language={language} copy={t} featured />
            ))}
          </div>
          <div className="grid">
            {recommendedEvents.map((event) => (
              <EventCard key={event.id} event={event} language={language} copy={t} />
            ))}
          </div>
        </section>

        <section id="cd-select" className="section">
          <div className="section-heading split-heading">
            <div>
              <p className="eyebrow">{t.cdEyebrow}</p>
              <h2>{t.cdTitle}</h2>
            </div>
            <p>{t.cdCopy}</p>
          </div>

          <div className="format-tabs" aria-label="CD select formats">
            {([
              ["all", t.formatAll],
              ["cd", t.formatCD],
              ["vinyl", t.formatVinyl],
            ] as const).map(([value, label]) => (
              <button
                key={value}
                className={formatFilter === value ? "active" : ""}
                type="button"
                onClick={() => setFormatFilter(value)}
              >
                {label}
              </button>
            ))}
          </div>

          <div className="catalog-grid">
            {filteredCatalog.map((item) => (
              <CatalogCard key={item.id} item={item} language={language} copy={t} />
            ))}
          </div>
        </section>

        <section id="archive" className="section">
          <div className="section-heading split-heading">
            <div>
              <p className="eyebrow">{t.archiveEyebrow}</p>
              <h2>{t.archiveTitle}</h2>
            </div>
            <p>{t.archiveCopy}</p>
          </div>

          <div className="grid compact-grid">
            {contents.map((item) => (
              <article className="content-card" key={item.id}>
                <span className="pill">{item.type}</span>
                <h3>{localized(item, "title", language)}</h3>
                <p>{localized(item, "summary", language)}</p>
              </article>
            ))}
          </div>
        </section>

        <section id="connect" className="section connect-section">
          <div className="connect-copy">
            <p className="eyebrow">{t.connectEyebrow}</p>
            <h2>{t.connectTitle}</h2>
            <p>{t.connectCopy}</p>
          </div>

          <form className="connect-form" onSubmit={submitConnect}>
            <label>
              <span>{t.labelNickname}</span>
              <input name="nickname" placeholder={t.placeholderNickname} required />
            </label>
            <label>
              <span>{t.labelEmail}</span>
              <input name="email" type="email" placeholder={t.placeholderEmail} required />
            </label>
            <label>
              <span>{t.labelTopic}</span>
              <select name="topic" defaultValue="ticket">
                <option value="ticket">{t.topicTicket}</option>
                <option value="cd-select">{t.topicCDSelect}</option>
                <option value="support">{t.topicSupport}</option>
                <option value="collab">{t.topicCollab}</option>
              </select>
            </label>
            <label className="wide">
              <span>{t.labelMessage}</span>
              <textarea name="message" rows={5} placeholder={t.placeholderMessage} />
            </label>
            <button type="submit">{t.submitButton}</button>
            <p className="form-result" role="status">{formStatus}</p>
          </form>
        </section>
      </main>
    </div>
  );
}

function V2Texture() {
  return (
    <>
      <V2BandScroll />
      <div className="band-texture" aria-hidden="true">
        THE SMITHS OASIS BLUR RADIOHEAD PULP SUEDE JOY DIVISION NEW ORDER THE CURE STONE ROSES MY BLOODY VALENTINE SLOWDIVE PRIMAL SCREAM MASSIVE ATTACK PORTISHEAD SONIC YOUTH PIXIES
      </div>
      <div className="track-field" aria-hidden="true">
        <span className="lane lane-a" />
        <span className="lane lane-b" />
        <span className="lane lane-c" />
        <span className="sample sample-a" />
        <span className="sample sample-b" />
        <span className="wave wave-a" />
        <span className="wave wave-b" />
      </div>
    </>
  );
}

function V2BandScroll() {
  const columns = [
    v2BandNames,
    [...v2BandNames].slice(28).concat(v2BandNames.slice(0, 28)),
    [...v2BandNames].slice(58).concat(v2BandNames.slice(0, 58)),
  ];

  return (
    <div className="v2-band-scroll" aria-hidden="true">
      {columns.map((names, index) => (
        <div className={`v2-band-column column-${index + 1}`} key={index}>
          {[...names, ...names].map((name, itemIndex) => (
            <span key={`${name}-${itemIndex}`}>{name}</span>
          ))}
        </div>
      ))}
    </div>
  );
}

function V3Signal() {
  return (
    <>
      <div className="v3-brand-signal" aria-hidden="true">
        <div className="v3-logo-mark" />
        <div className="v3-band-cloud">{classicBandCloud}</div>
        <div className="v3-marquee-row row-one">LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE</div>
        <div className="v3-marquee-row row-two">L I V E L I F E L I V E L I F E L I V E L I F E L I V E L I F E</div>
      </div>
      <div className="v3-track-field" aria-hidden="true">
        <div className="song-timeline">
          <span>THE LIBERTINES / MUSIC WHEN THE LIGHTS GO OUT</span>
          <span>126 BPM</span>
          <span>INDIE ROCK SESSION</span>
        </div>
        <span className="section-marker marker-a" />
        <span className="section-marker marker-b" />
        <span className="section-marker marker-c" />
        <span className="track-wave mix-wave-a" />
        <span className="track-row row-a" />
        <span className="track-row row-b" />
        <span className="track-row row-c" />
        <span className="track-row row-d" />
        <span className="track-block block-a" />
        <span className="track-block block-b" />
        <span className="track-block block-c" />
        <span className="beat-grid beat-a" />
        <span className="beat-grid beat-b" />
        <span className="track-clip clip-a" />
        <span className="track-clip clip-b" />
        <span className="track-clip clip-c" />
        <span className="track-wave mix-wave-b" />
        <span className="track-meter meter-a" />
        <span className="track-meter meter-b" />
      </div>
    </>
  );
}

function Entry({ href, index, label, icon }: { href: string; index: string; label: string; icon: ReactNode }) {
  return (
    <a className="entry-tile" href={href}>
      <span>
        <strong>{index}</strong>
        {icon}
      </span>
      <strong>{label}</strong>
    </a>
  );
}

function ScheduleRow({ item, language, copy }: { item: EventItem | CatalogItem; language: Language; copy: Copy }) {
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

function EventCard({ event, language, copy, featured = false }: { event: EventItem; language: Language; copy: Copy; featured?: boolean }) {
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

function CatalogCard({ item, language, copy }: { item: CatalogItem; language: Language; copy: Copy }) {
  return (
    <article className="catalog-card">
      {item.imageUrl ? (
        <img src={item.imageUrl} alt={localized(item, "title", language)} />
      ) : (
        <div className="catalog-art" aria-hidden="true" />
      )}
      <div className="catalog-body">
        <span className="pill">{item.format === "vinyl" ? copy.formatVinyl : copy.formatCD}</span>
        <h3>{localized(item, "title", language)}</h3>
        <p className="artist">{item.artist}</p>
        <p>{localized(item, "summary", language)}</p>
        <div className="track-list">{item.tracks?.map((track) => <span key={track}>{track}</span>)}</div>
        <div className="catalog-actions">
          <span>{item.price || "TBD"}</span>
          <a className="button buy-button" href={item.purchaseUrl} target="_blank" rel="noreferrer">
            {localizedMap(item.purchaseText, language) || "BUY HERE"}
          </a>
        </div>
        <small>{copy.externalShopNote}</small>
      </div>
    </article>
  );
}

function apiStatusLabel(status: "checking" | "online" | "offline", copy: Copy) {
  if (status === "online") return copy.apiOnline;
  if (status === "offline") return copy.apiOffline;
  return copy.apiChecking;
}

function localized<T extends { [key: string]: unknown }>(item: T, field: string, language: Language) {
  const i18n = item[`${field}I18n`] as LocalizedText | undefined;
  const fallback = item[field];
  return i18n?.[language] || i18n?.zh || (typeof fallback === "string" ? fallback : "");
}

function localizedMap(value: LocalizedText | undefined, language: Language) {
  return value?.[language] || value?.zh || "";
}
