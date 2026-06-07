import { Mail, Music2, ShoppingBag, Ticket, Vinyl } from "lucide-react";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

type Language = "zh" | "ja" | "en";

type LocalizedText = Partial<Record<Language, string>>;

type Event = {
  id: string;
  brand: string;
  owned: boolean;
  category: string;
  title: string;
  titleI18n?: LocalizedText;
  date: string;
  time: string;
  venue: string;
  area: string;
  price: string;
  tags: string[];
  summary: string;
  summaryI18n?: LocalizedText;
  lineup?: string[];
  ticketNote?: string;
  ticketNoteI18n?: LocalizedText;
  imageUrl?: string;
  sourceNote?: string;
  sourceNoteI18n?: LocalizedText;
};

type CatalogItem = {
  id: string;
  brand: string;
  kind: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
  status: string;
  price: string;
  imageUrl?: string;
};

type ContentItem = {
  id: string;
  brand: string;
  type: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
};

type Copy = {
  navShows: string;
  navCD: string;
  navShop: string;
  navConnect: string;
  heroEyebrow: string;
  heroTitle: string;
  heroLead: string;
  heroPrimary: string;
  apiChecking: string;
  apiOnline: string;
  apiOffline: string;
  showsEyebrow: string;
  showsTitle: string;
  showsNote: string;
  cdEyebrow: string;
  cdTitle: string;
  cdCopy: string;
  shopEyebrow: string;
  shopTitle: string;
  shopCopy: string;
  contentEyebrow: string;
  contentTitle: string;
  connectEyebrow: string;
  connectTitle: string;
  connectCopy: string;
  ownBadge: string;
  recommendBadge: string;
};

const copy: Record<Language, Copy> = {
  zh: {
    navShows: "演出情报",
    navCD: "CD",
    navShop: "Shop",
    navConnect: "Connect",
    heroEyebrow: "TOKYO MUSIC ENTRY",
    heroTitle: "LIVE LIFE 把东京现场、CD 和 Shop 分开整理。",
    heroLead: "先从我们自己的演出和推荐演出情报开始，把票务、CD、Shop 和售后联系拆成清楚的入口。",
    heroPrimary: "看演出情报",
    apiChecking: "API 检查中",
    apiOnline: "API 已连接",
    apiOffline: "API 未连接",
    showsEyebrow: "SHOWS",
    showsTitle: "演出情报",
    showsNote: "我们自己的演出固定在最上面；推荐演出和历史视觉档案放在下面。",
    cdEyebrow: "CD",
    cdTitle: "CD 单独成页",
    cdCopy: "CD 不和演出页混在一起，后续可以做试听、发行物介绍和购买入口。",
    shopEyebrow: "SHOP",
    shopTitle: "Shop 独立讨论购买流程",
    shopCopy: "登录注册、购物车、支付和订单先不做，等购买流程确定后再接。",
    contentEyebrow: "NOTES",
    contentTitle: "首页内容摘要",
    connectEyebrow: "CONNECT",
    connectTitle: "付款、票务、发货或合作问题，都从这里联系。",
    connectCopy: "这个入口不是 Join us，而是 LIVE LIFE 的统一消息入口。比如付款后没有收到货、CD/Shop 发货问题、活动合作、投稿，都可以从这里发消息。",
    ownBadge: "LIVE LIFE 自主演出",
    recommendBadge: "推荐 / 档案",
  },
  ja: {
    navShows: "ライブ情報",
    navCD: "CD",
    navShop: "Shop",
    navConnect: "Connect",
    heroEyebrow: "TOKYO MUSIC ENTRY",
    heroTitle: "LIVE LIFE は東京のライブ、CD、Shop を分けて整理します。",
    heroLead: "まずは自主公演とおすすめライブ情報から始め、チケット、CD、Shop、問い合わせの入口を明確に分けます。",
    heroPrimary: "ライブ情報を見る",
    apiChecking: "API 確認中",
    apiOnline: "API 接続済み",
    apiOffline: "API 未接続",
    showsEyebrow: "SHOWS",
    showsTitle: "ライブ情報",
    showsNote: "自主公演を最上部に固定し、おすすめ公演と過去ビジュアルは下に分けて表示します。",
    cdEyebrow: "CD",
    cdTitle: "CD は独立ページへ",
    cdCopy: "CD はライブ情報と混ぜず、試聴、リリース紹介、購入導線を後から追加できる形にします。",
    shopEyebrow: "SHOP",
    shopTitle: "Shop は購入フローを別途検討",
    shopCopy: "ログイン、カート、決済、注文管理はまだ作らず、購入フローが決まってから接続します。",
    contentEyebrow: "NOTES",
    contentTitle: "ホーム用メモ",
    connectEyebrow: "CONNECT",
    connectTitle: "支払い、チケット、発送、コラボの相談はこちらから。",
    connectCopy: "これは Join us ではなく、LIVE LIFE の共通問い合わせ入口です。未着、CD/Shop、イベント協力、投稿などをここから送れます。",
    ownBadge: "LIVE LIFE 自主公演",
    recommendBadge: "おすすめ / アーカイブ",
  },
  en: {
    navShows: "Shows",
    navCD: "CD",
    navShop: "Shop",
    navConnect: "Connect",
    heroEyebrow: "TOKYO MUSIC ENTRY",
    heroTitle: "LIVE LIFE separates Tokyo shows, CDs, and Shop clearly.",
    heroLead: "We start with owned shows and recommended live information, then keep ticketing, CDs, Shop, and support messages in separate lanes.",
    heroPrimary: "View shows",
    apiChecking: "Checking API",
    apiOnline: "API connected",
    apiOffline: "API offline",
    showsEyebrow: "SHOWS",
    showsTitle: "Live information",
    showsNote: "Owned LIVE LIFE shows stay at the top. Recommendations and archive visuals sit below.",
    cdEyebrow: "CD",
    cdTitle: "CD gets its own page",
    cdCopy: "CD content stays separate from shows, ready for listening notes, releases, and future purchase links.",
    shopEyebrow: "SHOP",
    shopTitle: "Shop stays separate while checkout is discussed",
    shopCopy: "Login, cart, payment, and orders are intentionally deferred until the buying flow is decided.",
    contentEyebrow: "NOTES",
    contentTitle: "Homepage notes",
    connectEyebrow: "CONNECT",
    connectTitle: "Payment, ticketing, shipping, or collaboration questions start here.",
    connectCopy: "This is not a Join us form. It is the shared LIVE LIFE message entry for missing orders, CD/Shop questions, event collaboration, and submissions.",
    ownBadge: "LIVE LIFE owned show",
    recommendBadge: "Recommendation / archive",
  },
};

export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [events, setEvents] = useState<Event[]>([]);
  const [cdItems, setCDItems] = useState<CatalogItem[]>([]);
  const [shopItems, setShopItems] = useState<CatalogItem[]>([]);
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [apiStatus, setApiStatus] = useState<"checking" | "online" | "offline">("checking");
  const t = copy[language];

  // React 前端正式启用时，也要同步 html lang。
  // 这样浏览器翻译、无障碍工具和搜索引擎能识别当前页面语言。
  useEffect(() => {
    document.documentElement.lang = language === "zh" ? "zh-Hans" : language;
  }, [language]);

  // 本地开发时 Vite 会把 /api 代理到 Go 后端。
  // 这里一次性拉取首页需要的演出、CD、Shop 和内容摘要，保持信息架构清楚分层。
  useEffect(() => {
    async function load() {
      try {
        const [healthRes, eventsRes, cdRes, shopRes, contentsRes] = await Promise.all([
          fetch("/api/health"),
          fetch("/api/events"),
          fetch("/api/cd-items"),
          fetch("/api/shop-items"),
          fetch("/api/contents"),
        ]);
        if (!healthRes.ok || !eventsRes.ok || !cdRes.ok || !shopRes.ok || !contentsRes.ok) {
          throw new Error("API request failed");
        }
        const eventsData = await eventsRes.json();
        const cdData = await cdRes.json();
        const shopData = await shopRes.json();
        const contentsData = await contentsRes.json();
        setEvents(eventsData.events || []);
        setCDItems(cdData.items || []);
        setShopItems(shopData.items || []);
        setContents(contentsData.contents || []);
        setApiStatus("online");
      } catch {
        setApiStatus("offline");
      }
    }

    load();
  }, []);

  const ownedEvents = events.filter((event) => event.owned);
  const recommendedEvents = events.filter((event) => !event.owned);

  return (
    <div className="min-h-screen bg-night text-cream">
      <header className="sticky top-0 z-10 border-b border-cream/10 bg-night/85 px-6 py-4 backdrop-blur">
        <div className="mx-auto flex max-w-6xl flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <a className="text-2xl font-extrabold" href="#home">LIVE LIFE</a>
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <nav className="flex flex-wrap gap-5 text-sm text-cream/65">
              <a href="#shows">{t.navShows}</a>
              <a href="#cd">{t.navCD}</a>
              <a href="#shop">{t.navShop}</a>
              <a href="#connect">{t.navConnect}</a>
            </nav>
            <LanguageSwitch language={language} onChange={setLanguage} />
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6">
        <section id="home" className="grid min-h-[82vh] items-center gap-10 py-16 md:grid-cols-[1fr_360px]">
          <div>
            <p className="text-sm font-extrabold uppercase text-amber">{t.heroEyebrow}</p>
            <h1 className="mt-3 max-w-4xl text-[44px] font-black leading-none md:text-[68px]">
              {t.heroTitle}
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-cream/70">{t.heroLead}</p>
            <a className="mt-8 inline-flex rounded-lg bg-cream px-5 py-3 font-extrabold text-night" href="#shows">
              {t.heroPrimary}
            </a>
          </div>
          <aside className="grid gap-4">
            <div className="rounded-lg border border-cream/10 bg-cream/5 p-4">
              <span className={`mb-4 block size-3 rounded-full ${apiStatus === "online" ? "bg-green" : "bg-amber"}`} />
              <strong>{apiStatusLabel(apiStatus, t)}</strong>
              <p className="mt-2 text-sm text-cream/60">LIVE LIFE API</p>
            </div>
            <img
              className="aspect-[4/5] rounded-lg object-cover shadow-2xl"
              src="/assets/events/redhair-2026-july.jpg"
              alt="紅髪少年殺人事件 2026 Tokyo poster"
            />
          </aside>
        </section>

        <SectionIntro id="shows" eyebrow={t.showsEyebrow} title={t.showsTitle} note={t.showsNote} icon={<Ticket size={20} />} />
        <div className="grid gap-5">
          {ownedEvents.map((event) => (
            <EventCard key={event.id} event={event} language={language} badge={t.ownBadge} featured />
          ))}
          <div className="grid gap-5 md:grid-cols-2">
            {recommendedEvents.map((event) => (
              <EventCard key={event.id} event={event} language={language} badge={t.recommendBadge} />
            ))}
          </div>
        </div>

        <CatalogSection id="cd" icon={<Vinyl size={20} />} eyebrow={t.cdEyebrow} title={t.cdTitle} copy={t.cdCopy} items={cdItems} language={language} />
        <CatalogSection id="shop" icon={<ShoppingBag size={20} />} eyebrow={t.shopEyebrow} title={t.shopTitle} copy={t.shopCopy} items={shopItems} language={language} />

        <section className="border-t border-cream/10 py-14">
          <p className="text-sm font-extrabold uppercase text-amber">{t.contentEyebrow}</p>
          <h2 className="mt-2 text-3xl font-black">{t.contentTitle}</h2>
          <div className="mt-6 grid gap-5 md:grid-cols-2">
            {contents.map((item) => (
              <article key={item.id} className="rounded-lg border border-cream/10 bg-cream/5 p-5">
                <span className="rounded-full border border-amber/40 px-3 py-1 text-xs font-extrabold text-amber">{item.type}</span>
                <h3 className="mt-4 text-2xl font-black">{localized(item, "title", language)}</h3>
                <p className="mt-3 leading-7 text-cream/65">{localized(item, "summary", language)}</p>
              </article>
            ))}
          </div>
        </section>

        <section id="connect" className="grid gap-7 border-t border-cream/10 py-14 md:grid-cols-[0.8fr_1.2fr]">
          <div>
            <p className="text-sm font-extrabold uppercase text-amber">{t.connectEyebrow}</p>
            <h2 className="mt-2 text-3xl font-black">{t.connectTitle}</h2>
            <p className="mt-4 leading-7 text-cream/65">{t.connectCopy}</p>
          </div>
          <div className="rounded-lg border border-cream/10 bg-cream/5 p-5">
            <div className="flex items-center gap-3 text-cream/75">
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
    <div className="inline-grid grid-cols-3 overflow-hidden rounded-lg border border-cream/10 bg-cream/5" aria-label="语言选择">
      {(["zh", "ja", "en"] as const).map((item) => (
        <button
          key={item}
          className={`min-w-[76px] border-l border-cream/10 px-3 py-2 text-sm font-bold first:border-l-0 ${
            language === item ? "bg-cream text-night" : "text-cream/65"
          }`}
          type="button"
          aria-pressed={language === item}
          onClick={() => onChange(item)}
        >
          {item === "zh" ? "中文" : item === "ja" ? "日本語" : "English"}
        </button>
      ))}
    </div>
  );
}

function SectionIntro({ id, eyebrow, title, note, icon }: { id: string; eyebrow: string; title: string; note: string; icon: ReactNode }) {
  return (
    <section id={id} className="border-t border-cream/10 pt-14">
      <div className="mb-6 grid gap-5 md:grid-cols-[1fr_360px] md:items-end">
        <div>
          <div className="mb-3 text-amber">{icon}</div>
          <p className="text-sm font-extrabold uppercase text-amber">{eyebrow}</p>
          <h2 className="mt-2 text-3xl font-black">{title}</h2>
        </div>
        <p className="leading-7 text-cream/65">{note}</p>
      </div>
    </section>
  );
}

function EventCard({ event, language, badge, featured = false }: { event: Event; language: Language; badge: string; featured?: boolean }) {
  return (
    <article className={`overflow-hidden rounded-lg border border-cream/10 bg-cream/5 ${featured ? "md:grid md:grid-cols-[410px_1fr]" : ""}`}>
      {event.imageUrl ? (
        <img className="aspect-[4/5] size-full object-cover" src={event.imageUrl} alt={localized(event, "title", language)} />
      ) : (
        <div className="aspect-[4/5] bg-cream/10" />
      )}
      <div className="p-5">
        <span className="rounded-full border border-amber/40 px-3 py-1 text-xs font-extrabold text-amber">{badge}</span>
        <h3 className="mt-4 text-2xl font-black">{localized(event, "title", language)}</h3>
        <div className="mt-4 flex flex-wrap gap-2 text-sm text-cream/70">
          {[event.date, event.time, event.venue, event.price].map((text) => (
            <span key={text} className="rounded-full bg-cream/10 px-3 py-1">{text}</span>
          ))}
        </div>
        <p className="mt-4 leading-7 text-cream/65">{localized(event, "summary", language)}</p>
      </div>
    </article>
  );
}

function CatalogSection({ id, icon, eyebrow, title, copy: sectionCopy, items, language }: { id: string; icon: ReactNode; eyebrow: string; title: string; copy: string; items: CatalogItem[]; language: Language }) {
  return (
    <section id={id} className="grid gap-7 border-t border-cream/10 py-14 md:grid-cols-[1fr_420px]">
      <div>
        <div className="mb-3 text-amber">{icon}</div>
        <p className="text-sm font-extrabold uppercase text-amber">{eyebrow}</p>
        <h2 className="mt-2 text-3xl font-black">{title}</h2>
        <p className="mt-4 leading-7 text-cream/65">{sectionCopy}</p>
      </div>
      <div className="grid gap-4">
        {items.map((item) => (
          <article key={item.id} className="rounded-lg border border-cream/10 bg-cream/5 p-5">
            <span className="rounded-full border border-amber/40 px-3 py-1 text-xs font-extrabold text-amber">{item.status}</span>
            <h3 className="mt-4 text-2xl font-black">{localized(item, "title", language)}</h3>
            <p className="mt-3 leading-7 text-cream/65">{localized(item, "summary", language)}</p>
            <strong>{item.price}</strong>
          </article>
        ))}
      </div>
    </section>
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
