import { Archive, Disc3, Mail, Ticket } from "lucide-react";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

type Language = "zh" | "ja" | "en";
type LocalizedText = Partial<Record<Language, string>>;

type Event = {
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
};

type ContentItem = {
  id: string;
  type: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
};

type Copy = {
  navShows: string;
  navCDSelect: string;
  navArchive: string;
  navConnect: string;
  heroTitle: string;
  heroLead: string;
  apiChecking: string;
  apiOnline: string;
  apiOffline: string;
  scheduleTitle: string;
  showsTitle: string;
  cdTitle: string;
  archiveTitle: string;
  connectTitle: string;
  buyNote: string;
  ownBadge: string;
};

const copy: Record<Language, Copy> = {
  zh: {
    navShows: "演出情报",
    navCDSelect: "CD 严选",
    navArchive: "档案",
    navConnect: "联系",
    heroTitle: "LIVE LIFE 是东京现场、CD 严选和音乐档案入口。",
    heroLead: "我们把自主演出、推荐现场、CD/黑胶严选、历史档案和售后联系整理成一个清楚的音乐入口。",
    apiChecking: "API 检查中",
    apiOnline: "API 已连接",
    apiOffline: "API 未连接",
    scheduleTitle: "近期日程",
    showsTitle: "演出情报",
    cdTitle: "CD 严选",
    archiveTitle: "档案",
    connectTitle: "票务、购买、发货或合作问题，都从这里联系。",
    buyNote: "购买会跳转到外部 Shop",
    ownBadge: "LIVE LIFE 自主演出",
  },
  ja: {
    navShows: "ライブ情報",
    navCDSelect: "CDセレクト",
    navArchive: "アーカイブ",
    navConnect: "問い合わせ",
    heroTitle: "LIVE LIFE は東京のライブ、CDセレクト、音楽アーカイブの入口です。",
    heroLead: "自主公演、おすすめライブ、CD/ヴァイナルセレクト、アーカイブ、問い合わせをひとつの音楽入口として整理します。",
    apiChecking: "API 確認中",
    apiOnline: "API 接続済み",
    apiOffline: "API 未接続",
    scheduleTitle: "近日予定",
    showsTitle: "ライブ情報",
    cdTitle: "CDセレクト",
    archiveTitle: "アーカイブ",
    connectTitle: "チケット、購入、発送、コラボの相談はこちらから。",
    buyNote: "購入は外部ショップへ移動します",
    ownBadge: "LIVE LIFE 自主公演",
  },
  en: {
    navShows: "SHOWS",
    navCDSelect: "CD SELECT",
    navArchive: "ARCHIVE",
    navConnect: "CONNECT",
    heroTitle: "LIVE LIFE IS AN ENTRY POINT FOR TOKYO SHOWS, CD SELECT, AND MUSIC ARCHIVES.",
    heroLead: "WE ORGANIZE OWNED SHOWS, RECOMMENDED LIVE DATES, CD/VINYL SELECT, ARCHIVES, AND SUPPORT MESSAGES INTO ONE MUSIC INDEX.",
    apiChecking: "CHECKING API",
    apiOnline: "API CONNECTED",
    apiOffline: "API OFFLINE",
    scheduleTitle: "UPCOMING",
    showsTitle: "SHOWS",
    cdTitle: "CD SELECT",
    archiveTitle: "ARCHIVE",
    connectTitle: "TICKETS, PURCHASES, SHIPPING, OR COLLABORATION QUESTIONS START HERE.",
    buyNote: "PURCHASE OPENS AN EXTERNAL SHOP",
    ownBadge: "LIVE LIFE OWNED SHOW",
  },
};

export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [events, setEvents] = useState<Event[]>([]);
  const [catalog, setCatalog] = useState<CatalogItem[]>([]);
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [apiStatus, setApiStatus] = useState<"checking" | "online" | "offline">("checking");
  const t = copy[language];

  useEffect(() => {
    document.documentElement.lang = language === "zh" ? "zh-Hans" : language;
  }, [language]);

  useEffect(() => {
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
        const eventsData = await eventsRes.json();
        const cdData = await cdRes.json();
        const contentsData = await contentsRes.json();
        setEvents(eventsData.events || []);
        setCatalog(cdData.items || []);
        setContents(contentsData.contents || []);
        setApiStatus("online");
      } catch {
        setApiStatus("offline");
      }
    }

    load();
  }, []);

  return (
    <div className="min-h-screen bg-paper text-ink">
      <header className="sticky top-0 z-10 border-b border-red bg-paper/90 px-6 py-4 backdrop-blur">
        <div className="mx-auto flex max-w-6xl flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <a className="flex items-center gap-3 font-black" href="#home">
            <span className="size-8 border-2 border-ink bg-yellow shadow-[5px_5px_0_#2457ff]" />
            <span>LIVE LIFE</span>
          </a>
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <nav className="flex flex-wrap gap-5 text-sm font-black">
              <a href="#shows">{t.navShows}</a>
              <a href="#cd-select">{t.navCDSelect}</a>
              <a href="#archive">{t.navArchive}</a>
              <a href="#connect">{t.navConnect}</a>
            </nav>
            <LanguageSwitch language={language} onChange={setLanguage} />
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6">
        <section id="home" className="grid min-h-[78vh] gap-0 border-x border-ink/15 md:grid-cols-[1fr_360px]">
          <div className="relative min-h-[620px] overflow-hidden border-r border-ink/15">
            <div className="absolute inset-0 p-8 text-[48px] font-black uppercase leading-none text-ink/15 md:text-[72px]">
              THE SMITHS / OASIS / BLUR / RADIOHEAD / PULP / SUEDE / JOY DIVISION / NEW ORDER / THE CURE / SLOWDIVE / PIXIES
            </div>
            <div className="absolute left-1/2 top-1/2 w-[min(620px,calc(100%-40px))] -translate-x-1/2 -translate-y-1/2 border-4 border-ink bg-paper p-8 shadow-[12px_12px_0_#ffd000,-9px_-9px_0_#2457ff]">
              <p className="text-sm font-black text-red">TOKYO MUSIC INDEX</p>
              <h1 className="mt-3 text-5xl font-black leading-none md:text-7xl">{t.heroTitle}</h1>
              <p className="mt-5 text-lg leading-8 text-muted">{t.heroLead}</p>
            </div>
          </div>
          <aside className="bg-red p-6 text-white">
            <div className="mb-8 border-b border-white/40 pb-4">
              <span className={`mb-3 block size-3 rounded-full ${apiStatus === "online" ? "bg-green" : "bg-yellow"}`} />
              <strong>{apiStatusLabel(apiStatus, t)}</strong>
            </div>
            <h2 className="text-4xl font-black">{t.scheduleTitle}</h2>
            <div className="mt-5 grid">
              {[...events.filter((event) => event.owned).slice(0, 2), ...catalog.slice(0, 2)].map((item) => (
                <ScheduleRow key={item.id} item={item} language={language} />
              ))}
            </div>
          </aside>
        </section>

        <section className="grid border-x border-b border-ink/15 md:grid-cols-4">
          <Entry href="#shows" index="01" label={t.navShows} icon={<Ticket />} />
          <Entry href="#cd-select" index="02" label={t.navCDSelect} icon={<Disc3 />} />
          <Entry href="#archive" index="03" label={t.navArchive} icon={<Archive />} />
          <Entry href="#connect" index="04" label={t.navConnect} icon={<Mail />} />
        </section>

        <section id="shows" className="border-t border-red py-14">
          <h2 className="text-5xl font-black">{t.showsTitle}</h2>
          <div className="mt-6 grid gap-5">
            {events.map((event) => (
              <article key={event.id} className="grid overflow-hidden border border-ink/15 bg-white/40 md:grid-cols-[320px_1fr]">
                {event.imageUrl && <img className="aspect-[4/5] h-full w-full object-cover" src={event.imageUrl} alt={localized(event, "title", language)} />}
                <div className="p-5">
                  {event.owned && <span className="border border-blue px-2 py-1 text-xs font-black text-blue">{t.ownBadge}</span>}
                  <h3 className="mt-4 text-3xl font-black">{localized(event, "title", language)}</h3>
                  <p className="mt-3 text-muted">{event.date} / {event.venue}</p>
                  <p className="mt-4 leading-7 text-muted">{localized(event, "summary", language)}</p>
                </div>
              </article>
            ))}
          </div>
        </section>

        <section id="cd-select" className="border-t border-red py-14">
          <h2 className="text-5xl font-black">{t.cdTitle}</h2>
          <div className="mt-6 grid gap-5 md:grid-cols-2">
            {catalog.map((item) => (
              <article key={item.id} className="grid overflow-hidden border border-ink/15 bg-white/40 md:grid-cols-[170px_1fr]">
                {item.imageUrl ? <img className="aspect-[4/5] h-full w-full object-cover" src={item.imageUrl} alt={localized(item, "title", language)} /> : <div className="min-h-[240px] bg-yellow" />}
                <div className="p-5">
                  <span className="border border-blue px-2 py-1 text-xs font-black uppercase text-blue">{item.format}</span>
                  <h3 className="mt-4 text-2xl font-black">{localized(item, "title", language)}</h3>
                  <p className="font-black text-red">{item.artist}</p>
                  <p className="mt-3 leading-7 text-muted">{localized(item, "summary", language)}</p>
                  <a className="mt-5 inline-flex bg-red px-4 py-2 font-black text-white" href={item.purchaseUrl} target="_blank" rel="noreferrer">
                    {item.purchaseText?.[language] || item.purchaseText?.zh || "BUY HERE"}
                  </a>
                  <small className="mt-2 block text-muted">{t.buyNote}</small>
                </div>
              </article>
            ))}
          </div>
        </section>

        <section id="archive" className="border-t border-red py-14">
          <h2 className="text-5xl font-black">{t.archiveTitle}</h2>
          <div className="mt-6 grid gap-5 md:grid-cols-2">
            {contents.map((item) => (
              <article key={item.id} className="border border-ink/15 bg-white/40 p-5">
                <span className="border border-blue px-2 py-1 text-xs font-black uppercase text-blue">{item.type}</span>
                <h3 className="mt-4 text-2xl font-black">{localized(item, "title", language)}</h3>
                <p className="mt-3 leading-7 text-muted">{localized(item, "summary", language)}</p>
              </article>
            ))}
          </div>
        </section>

        <section id="connect" className="grid gap-7 border-t border-red py-14 md:grid-cols-[0.8fr_1.2fr]">
          <div>
            <p className="text-sm font-black text-red">CONNECT</p>
            <h2 className="mt-2 text-4xl font-black">{t.connectTitle}</h2>
          </div>
          <div className="border border-ink/15 bg-white/40 p-5">
            <div className="flex items-center gap-3 text-muted">
              <Mail size={20} />
              <span>/api/connect</span>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}

function LanguageSwitch({ language, onChange }: { language: Language; onChange: (language: Language) => void }) {
  return (
    <div className="inline-grid grid-cols-3 border border-ink">
      {(["zh", "ja", "en"] as const).map((item) => (
        <button
          key={item}
          className={`min-w-[76px] border-l border-ink px-3 py-2 text-sm font-black first:border-l-0 ${language === item ? "bg-ink text-paper" : ""}`}
          type="button"
          aria-pressed={language === item}
          onClick={() => onChange(item)}
        >
          {item === "zh" ? "中文" : item === "ja" ? "日本語" : "ENGLISH"}
        </button>
      ))}
    </div>
  );
}

function Entry({ href, index, label, icon }: { href: string; index: string; label: string; icon: ReactNode }) {
  return (
    <a className="flex min-h-40 flex-col justify-between border-l border-ink/15 p-5 first:border-l-0" href={href}>
      <span className="flex items-center justify-between text-blue">
        <strong>{index}</strong>
        {icon}
      </span>
      <strong className="text-4xl font-black uppercase leading-none text-red">{label}</strong>
    </a>
  );
}

function ScheduleRow({ item, language }: { item: Event | CatalogItem; language: Language }) {
  const isEvent = "date" in item;
  return (
    <article className="grid grid-cols-[86px_1fr] gap-4 border-t border-white/40 py-4">
      <strong>{isEvent ? item.date : "TBD"}</strong>
      <div>
        <span className="text-xs font-black text-white/70">{isEvent ? "LIVE" : item.format.toUpperCase()}</span>
        <p className="font-black">{localized(item, "title", language)}</p>
        <small className="text-white/70">{isEvent ? item.venue : item.artist}</small>
      </div>
    </article>
  );
}

function apiStatusLabel(status: "checking" | "online" | "offline", t: Copy) {
  if (status === "online") return t.apiOnline;
  if (status === "offline") return t.apiOffline;
  return t.apiChecking;
}

function localized<T extends { title: string; summary: string }>(item: T & { titleI18n?: LocalizedText; summaryI18n?: LocalizedText }, field: "title" | "summary", language: Language) {
  const i18n = field === "title" ? item.titleI18n : item.summaryI18n;
  return i18n?.[language] || i18n?.zh || item[field];
}
