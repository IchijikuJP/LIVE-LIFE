import { Archive, Disc3, Mail, Ticket } from "lucide-react";
import { FormEvent, ReactNode, useEffect, useMemo, useState } from "react";

type Language = "zh" | "ja" | "en";
type DesignVariant = "v2" | "v3";
type FormatFilter = "all" | "cd" | "vinyl";
type LocalizedText = Partial<Record<Language, string>>;

type EventItem = {
  id: string;
  owned: boolean;
  title: string;
  titleI18n?: LocalizedText;
  date: string;
  time: string;
  venue: string;
  price: string;
  summary: string;
  summaryI18n?: LocalizedText;
  imageUrl?: string;
  lineup?: string[];
  tags?: string[];
  ticketNote?: string;
  ticketNoteI18n?: LocalizedText;
  sourceNote?: string;
  sourceNoteI18n?: LocalizedText;
};

type CatalogItem = {
  id: string;
  format: "cd" | "vinyl";
  artist: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
  status: string;
  price: string;
  imageUrl?: string;
  purchaseUrl: string;
  purchaseText?: LocalizedText;
  tracks?: string[];
};

type ContentItem = {
  id: string;
  type: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
};

type HealthPayload = {
  status: string;
  brand: string;
  service: string;
};

type Copy = {
  langAttr: string;
  navShows: string;
  navCDSelect: string;
  navArchive: string;
  navConnect: string;
  designLabel: string;
  designV2: string;
  designV3: string;
  heroEyebrow: string;
  heroTitle: string;
  heroLead: string;
  heroPrimary: string;
  heroSecondary: string;
  apiChecking: string;
  apiOnline: string;
  apiOffline: string;
  scheduleEyebrow: string;
  scheduleTitle: string;
  showsEyebrow: string;
  showsTitle: string;
  showsNote: string;
  cdEyebrow: string;
  cdTitle: string;
  cdCopy: string;
  formatAll: string;
  formatCD: string;
  formatVinyl: string;
  archiveEyebrow: string;
  archiveTitle: string;
  archiveCopy: string;
  connectEyebrow: string;
  connectTitle: string;
  connectCopy: string;
  labelNickname: string;
  labelEmail: string;
  labelTopic: string;
  labelMessage: string;
  placeholderNickname: string;
  placeholderEmail: string;
  placeholderMessage: string;
  topicTicket: string;
  topicCDSelect: string;
  topicSupport: string;
  topicCollab: string;
  submitButton: string;
  submitting: string;
  submitFallback: string;
  submitFailed: string;
  ownBadge: string;
  recommendBadge: string;
  ticketPending: string;
  sourceNote: string;
  scheduleLive: string;
  scheduleCD: string;
  scheduleVinyl: string;
  externalShopNote: string;
};

const copy: Record<Language, Copy> = {
  zh: {
    langAttr: "zh-Hans",
    navShows: "演出情报",
    navCDSelect: "CD 严选",
    navArchive: "档案",
    navConnect: "联系",
    designLabel: "设计方案",
    designV2: "V2 当前版",
    designV3: "V3 乐队信号版",
    heroEyebrow: "TOKYO MUSIC INDEX",
    heroTitle: "LIVE LIFE 是东京现场、CD 严选和音乐档案入口。",
    heroLead: "我们把 LIVE LIFE 自主演出、推荐现场、CD/黑胶严选、历史档案和售后联系整理成一个清晰的音乐入口。",
    heroPrimary: "看演出日程",
    heroSecondary: "进入 CD 严选",
    apiChecking: "API 检查中",
    apiOnline: "API 已连接",
    apiOffline: "API 未连接",
    scheduleEyebrow: "SCHEDULE",
    scheduleTitle: "近期日程",
    showsEyebrow: "SHOWS",
    showsTitle: "演出情报",
    showsNote: "LIVE LIFE 自主演出固定在最上方，推荐演出和历史视觉档案放在下面。",
    cdEyebrow: "CD SELECT",
    cdTitle: "CD 严选",
    cdCopy: "这里分成 CD 和黑胶。单品详情里的购买按钮会跳到外部 Shop，例如 BASE。",
    formatAll: "全部",
    formatCD: "CD",
    formatVinyl: "黑胶",
    archiveEyebrow: "ARCHIVE",
    archiveTitle: "档案",
    archiveCopy: "历史海报、公开资料备注、照片和推荐文章以后会集中在这里。",
    connectEyebrow: "CONNECT",
    connectTitle: "票务、购买、发货或合作问题，都从这里联系。",
    connectCopy: "外部平台购买后没有收到货、CD/黑胶购买咨询、票务问题、活动合作和投稿，都可以从这里给 LIVE LIFE 发消息。",
    labelNickname: "昵称",
    labelEmail: "邮箱",
    labelTopic: "问题类型",
    labelMessage: "留言",
    placeholderNickname: "LIVE LIFE 朋友",
    placeholderEmail: "you@example.com",
    placeholderMessage: "请写下你遇到的问题或想联系 LIVE LIFE 的原因。",
    topicTicket: "票务",
    topicCDSelect: "CD 严选",
    topicSupport: "购买或发货问题",
    topicCollab: "合作 / 投稿",
    submitButton: "发送消息",
    submitting: "发送中...",
    submitFallback: "LIVE LIFE 已收到你的消息。",
    submitFailed: "发送失败",
    ownBadge: "LIVE LIFE 自主演出",
    recommendBadge: "推荐 / 档案",
    ticketPending: "票务待确认",
    sourceNote: "资料备注",
    scheduleLive: "LIVE",
    scheduleCD: "CD",
    scheduleVinyl: "VINYL",
    externalShopNote: "购买会跳转到外部 Shop",
  },
  ja: {
    langAttr: "ja",
    navShows: "ライブ情報",
    navCDSelect: "CD セレクト",
    navArchive: "アーカイブ",
    navConnect: "問い合わせ",
    designLabel: "デザイン",
    designV2: "V2 現行案",
    designV3: "V3 バンド信号案",
    heroEyebrow: "TOKYO MUSIC INDEX",
    heroTitle: "LIVE LIFE は東京のライブ、CD セレクト、音楽アーカイブの入口です。",
    heroLead: "LIVE LIFE の自主公演、おすすめライブ、CD/ヴァイナルのセレクト、過去資料、問い合わせをひとつの音楽入口として整理します。",
    heroPrimary: "スケジュールを見る",
    heroSecondary: "CD セレクトへ",
    apiChecking: "API 確認中",
    apiOnline: "API 接続済み",
    apiOffline: "API 未接続",
    scheduleEyebrow: "SCHEDULE",
    scheduleTitle: "近日予定",
    showsEyebrow: "SHOWS",
    showsTitle: "ライブ情報",
    showsNote: "LIVE LIFE の自主公演を最上部に固定し、おすすめ公演と過去ビジュアルを下に配置します。",
    cdEyebrow: "CD SELECT",
    cdTitle: "CD セレクト",
    cdCopy: "CD とヴァイナルに分けます。詳細内の購入ボタンは BASE など外部 Shop へ移動します。",
    formatAll: "すべて",
    formatCD: "CD",
    formatVinyl: "ヴァイナル",
    archiveEyebrow: "ARCHIVE",
    archiveTitle: "アーカイブ",
    archiveCopy: "過去フライヤー、公開資料メモ、写真、推薦記事などをここに集約します。",
    connectEyebrow: "CONNECT",
    connectTitle: "チケット、購入、発送、コラボの相談はこちらから。",
    connectCopy: "外部 Shop 購入後の未着、CD/ヴァイナル購入相談、チケット、イベント協力、投稿などを LIVE LIFE に送れます。",
    labelNickname: "ニックネーム",
    labelEmail: "メール",
    labelTopic: "問い合わせ種別",
    labelMessage: "メッセージ",
    placeholderNickname: "LIVE LIFE の友人",
    placeholderEmail: "you@example.com",
    placeholderMessage: "困っていること、または LIVE LIFE に連絡したい内容を書いてください。",
    topicTicket: "チケット",
    topicCDSelect: "CD セレクト",
    topicSupport: "購入・発送",
    topicCollab: "協力 / 投稿",
    submitButton: "送信",
    submitting: "送信中...",
    submitFallback: "LIVE LIFE がメッセージを受け取りました。",
    submitFailed: "送信に失敗しました",
    ownBadge: "LIVE LIFE 自主公演",
    recommendBadge: "おすすめ / アーカイブ",
    ticketPending: "チケット確認中",
    sourceNote: "情報メモ",
    scheduleLive: "LIVE",
    scheduleCD: "CD",
    scheduleVinyl: "VINYL",
    externalShopNote: "購入は外部 Shop へ移動します",
  },
  en: {
    langAttr: "en",
    navShows: "SHOWS",
    navCDSelect: "CD SELECT",
    navArchive: "ARCHIVE",
    navConnect: "CONNECT",
    designLabel: "DESIGN",
    designV2: "V2 CURRENT",
    designV3: "V3 BAND SIGNAL",
    heroEyebrow: "TOKYO MUSIC INDEX",
    heroTitle: "LIVE LIFE IS AN ENTRY POINT FOR TOKYO SHOWS, CD SELECT, AND MUSIC ARCHIVES.",
    heroLead: "WE ORGANIZE LIVE LIFE OWNED SHOWS, RECOMMENDED LIVE DATES, CD/VINYL SELECT, ARCHIVES, AND SUPPORT MESSAGES INTO ONE CLEAR MUSIC INDEX.",
    heroPrimary: "VIEW SCHEDULE",
    heroSecondary: "ENTER CD SELECT",
    apiChecking: "CHECKING API",
    apiOnline: "API CONNECTED",
    apiOffline: "API OFFLINE",
    scheduleEyebrow: "SCHEDULE",
    scheduleTitle: "UPCOMING",
    showsEyebrow: "SHOWS",
    showsTitle: "SHOWS",
    showsNote: "LIVE LIFE OWNED SHOWS STAY AT THE TOP. RECOMMENDATIONS AND ARCHIVE VISUALS SIT BELOW.",
    cdEyebrow: "CD SELECT",
    cdTitle: "CD SELECT",
    cdCopy: "CD SELECT IS SPLIT INTO CD AND VINYL. PURCHASE BUTTONS ON DETAIL CARDS LINK TO EXTERNAL SHOPS SUCH AS BASE.",
    formatAll: "ALL",
    formatCD: "CD",
    formatVinyl: "VINYL",
    archiveEyebrow: "ARCHIVE",
    archiveTitle: "ARCHIVE",
    archiveCopy: "PAST POSTERS, PUBLIC RESEARCH NOTES, PHOTOS, AND RECOMMENDED ARTICLES WILL LIVE HERE.",
    connectEyebrow: "CONNECT",
    connectTitle: "TICKETS, PURCHASES, SHIPPING, OR COLLABORATION QUESTIONS START HERE.",
    connectCopy: "MESSAGE LIVE LIFE ABOUT EXTERNAL SHOP ORDERS, CD/VINYL QUESTIONS, TICKETS, EVENT COLLABORATION, OR SUBMISSIONS.",
    labelNickname: "NAME",
    labelEmail: "EMAIL",
    labelTopic: "TOPIC",
    labelMessage: "MESSAGE",
    placeholderNickname: "LIVE LIFE FRIEND",
    placeholderEmail: "you@example.com",
    placeholderMessage: "TELL US WHAT HAPPENED OR WHY YOU WANT TO CONTACT LIVE LIFE.",
    topicTicket: "TICKETING",
    topicCDSelect: "CD SELECT",
    topicSupport: "PURCHASE OR SHIPPING",
    topicCollab: "COLLABORATION / SUBMISSION",
    submitButton: "SEND MESSAGE",
    submitting: "SENDING...",
    submitFallback: "LIVE LIFE RECEIVED YOUR MESSAGE.",
    submitFailed: "SEND FAILED",
    ownBadge: "LIVE LIFE OWNED SHOW",
    recommendBadge: "RECOMMENDATION / ARCHIVE",
    ticketPending: "TICKETING PENDING",
    sourceNote: "SOURCE NOTE",
    scheduleLive: "LIVE",
    scheduleCD: "CD",
    scheduleVinyl: "VINYL",
    externalShopNote: "PURCHASE OPENS AN EXTERNAL SHOP",
  },
};

const classicBandCloud =
  "THE BEATLES THE ROLLING STONES LED ZEPPELIN PINK FLOYD QUEEN THE WHO DAVID BOWIE THE KINKS T. REX ROXY MUSIC THE CLASH SEX PISTOLS THE JAM BUZZCOCKS JOY DIVISION NEW ORDER THE SMITHS THE CURE SIOUXSIE AND THE BANSHEES ECHO AND THE BUNNYMEN THE STONE ROSES MY BLOODY VALENTINE SLOWDIVE PRIMAL SCREAM RADIOHEAD OASIS BLUR SUEDE PULP THE VERVE ARCTIC MONKEYS NIRVANA PIXIES SONIC YOUTH R.E.M. THE STOOGES THE VELVET UNDERGROUND TALKING HEADS";

const languageLabels: Record<Language, string> = {
  zh: "中文",
  ja: "日本語",
  en: "ENGLISH",
};

function readInitialDesign(): DesignVariant {
  const params = new URLSearchParams(window.location.search);
  const fromUrl = params.get("design");
  if (fromUrl === "v2" || fromUrl === "v3") {
    return fromUrl;
  }
  return window.localStorage.getItem("liveLifeDesignVariant") === "v3" ? "v3" : "v2";
}

export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [design, setDesign] = useState<DesignVariant>(readInitialDesign);
  const [formatFilter, setFormatFilter] = useState<FormatFilter>("all");
  const [events, setEvents] = useState<EventItem[]>([]);
  const [catalog, setCatalog] = useState<CatalogItem[]>([]);
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [health, setHealth] = useState<HealthPayload | null>(null);
  const [apiStatus, setApiStatus] = useState<"checking" | "online" | "offline">("checking");
  const [formStatus, setFormStatus] = useState("");
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
    url.searchParams.set("v", "20260608-db-v3");
    window.history.replaceState({}, "", url);
  }, [design]);

  // 前端启动后统一从 Go API 读取数据；这些 API 背后现在已经由 GORM 从 SQLite 读写。
  // 后续如果换成 PostgreSQL 或后台管理系统，React 页面不需要改字段，只要 API 契约保持一致。
  useEffect(() => {
    let alive = true;

    async function load() {
      try {
        const [healthRes, eventsRes, cdRes, contentsRes] = await Promise.all([
          fetch("/api/health"),
          fetch("/api/events"),
          fetch("/api/cd-items"),
          fetch("/api/contents"),
        ]);
        if (!healthRes.ok || !eventsRes.ok || !cdRes.ok || !contentsRes.ok) {
          throw new Error("API request failed");
        }
        const [healthData, eventsData, cdData, contentsData] = await Promise.all([
          healthRes.json(),
          eventsRes.json(),
          cdRes.json(),
          contentsRes.json(),
        ]);
        if (!alive) return;
        setHealth(healthData);
        setEvents(eventsData.events || []);
        setCatalog(cdData.items || []);
        setContents(contentsData.contents || []);
        setApiStatus("online");
      } catch {
        if (!alive) return;
        setApiStatus("offline");
      }
    }

    load();
    return () => {
      alive = false;
    };
  }, []);

  async function submitConnect(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormStatus(t.submitting);

    const form = event.currentTarget;
    const payload = Object.fromEntries(new FormData(form).entries());

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

function V3Signal() {
  return (
    <>
      <div className="v3-brand-signal" aria-hidden="true">
        <div className="v3-band-cloud">{classicBandCloud}</div>
        <div className="v3-marquee-row row-one">LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE</div>
        <div className="v3-marquee-row row-two">L I V E L I F E L I V E L I F E L I V E L I F E L I V E L I F E</div>
      </div>
      <div className="v3-track-field" aria-hidden="true">
        <span className="track-row row-a" />
        <span className="track-row row-b" />
        <span className="track-row row-c" />
        <span className="track-row row-d" />
        <span className="track-block block-a" />
        <span className="track-block block-b" />
        <span className="track-block block-c" />
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
