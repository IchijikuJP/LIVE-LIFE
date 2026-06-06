import { Camera, FileText, Send, Share2, ShoppingBag } from "lucide-react";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

type Language = "zh" | "ja";

type Event = {
  id: string;
  title: string;
  titleI18n?: Partial<Record<Language, string>>;
  date: string;
  time: string;
  venue: string;
  area: string;
  price: string;
  tags: string[];
  summary: string;
  summaryI18n?: Partial<Record<Language, string>>;
};

type ContentItem = {
  id: string;
  type: string;
  title: string;
  titleI18n?: Partial<Record<Language, string>>;
  summary: string;
  summaryI18n?: Partial<Record<Language, string>>;
};

type Copy = {
  navEvents: string;
  navContent: string;
  navConnect: string;
  navJoin: string;
  heroEyebrow: string;
  heroTitle: string;
  heroLead: string;
  apiChecking: string;
  apiOnline: string;
  apiOffline: string;
  proxyTarget: string;
  eventsEyebrow: string;
  eventsTitle: string;
  contentEyebrow: string;
  contentTitle: string;
  connectEyebrow: string;
  connectTitle: string;
  linkArticles: string;
  linkPhotos: string;
  linkShop: string;
  linkSubmission: string;
  joinEyebrow: string;
  joinTitle: string;
  joinNoteBefore: string;
  joinNoteAfter: string;
};

const copy: Record<Language, Copy> = {
  zh: {
    navEvents: "活动 / 商店",
    navContent: "内容推荐",
    navConnect: "连接",
    navJoin: "加入我们",
    heroEyebrow: "东京 Livehouse MVP",
    heroTitle: "演出情报、CD/唱片 Shop 和本地音乐现场入口。",
    heroLead: "这是 LiveLife MVP 的 React 前端骨架。Go API 会在本地测试时提供活动、内容和 Join 表单接口。",
    apiChecking: "API 检查中",
    apiOnline: "API 已连接",
    apiOffline: "API 未连接",
    proxyTarget: "代理目标",
    eventsEyebrow: "下一场活动",
    eventsTitle: "活动 / 商店列表",
    contentEyebrow: "内容",
    contentTitle: "推荐内容",
    connectEyebrow: "连接",
    connectTitle: "链接入口",
    linkArticles: "文章",
    linkPhotos: "活动照片",
    linkShop: "商店",
    linkSubmission: "投稿入口",
    joinEyebrow: "加入我们",
    joinTitle: "Join 表单占位",
    joinNoteBefore: "本地静态预览页已经可以提交到",
    joinNoteAfter: "。React 表单会复用同一个接口。",
  },
  ja: {
    navEvents: "イベント / ショップ",
    navContent: "おすすめ",
    navConnect: "つながる",
    navJoin: "参加する",
    heroEyebrow: "東京ライブハウス MVP",
    heroTitle: "ライブ情報、CD/レコードショップ、ローカルシーンへの入口。",
    heroLead: "LiveLife MVP の React フロントエンド骨組みです。ローカルテストでは Go API がイベント、コンテンツ、Join フォームを提供します。",
    apiChecking: "API 確認中",
    apiOnline: "API 接続済み",
    apiOffline: "API 未接続",
    proxyTarget: "プロキシ先",
    eventsEyebrow: "次のイベント",
    eventsTitle: "イベント / ショップ一覧",
    contentEyebrow: "コンテンツ",
    contentTitle: "おすすめ",
    connectEyebrow: "つながる",
    connectTitle: "リンク入口",
    linkArticles: "記事",
    linkPhotos: "イベント写真",
    linkShop: "ショップ",
    linkSubmission: "投稿入口",
    joinEyebrow: "参加する",
    joinTitle: "Join フォームのプレースホルダー",
    joinNoteBefore: "ローカル静的プレビューはすでに",
    joinNoteAfter: "へ送信できます。React フォームも同じエンドポイントを使います。",
  },
};

export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [events, setEvents] = useState<Event[]>([]);
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [apiStatus, setApiStatus] = useState<"checking" | "online" | "offline">("checking");
  const t = copy[language];

  useEffect(() => {
    document.documentElement.lang = language === "zh" ? "zh-Hans" : "ja";
  }, [language]);

  useEffect(() => {
    async function load() {
      try {
        const [healthRes, eventsRes, contentsRes] = await Promise.all([
          fetch("/api/health"),
          fetch("/api/events"),
          fetch("/api/contents"),
        ]);
        if (!healthRes.ok || !eventsRes.ok || !contentsRes.ok) {
          throw new Error("API request failed");
        }
        const eventsData = await eventsRes.json();
        const contentsData = await contentsRes.json();
        setEvents(eventsData.events);
        setContents(contentsData.contents);
        setApiStatus("online");
      } catch {
        setApiStatus("offline");
      }
    }

    load();
  }, []);

  return (
    <div className="min-h-screen bg-paper text-ink">
      <header className="sticky top-0 z-10 border-b border-stone-200 bg-paper/90 px-6 py-4 backdrop-blur">
        <div className="mx-auto flex max-w-6xl flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <strong className="text-2xl">LiveLife</strong>
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            <nav className="flex flex-wrap gap-5 text-sm text-stone-600">
              <a href="#events">{t.navEvents}</a>
              <a href="#content">{t.navContent}</a>
              <a href="#connect">{t.navConnect}</a>
              <a href="#join">{t.navJoin}</a>
            </nav>
            <LanguageSwitch language={language} onChange={setLanguage} />
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6">
        <section className="grid min-h-[72vh] items-center gap-8 py-14 md:grid-cols-[1fr_320px]">
          <div>
            <p className="text-sm font-bold uppercase text-moss">{t.heroEyebrow}</p>
            <h1 className="mt-3 max-w-4xl text-5xl font-black leading-none md:text-7xl">
              {t.heroTitle}
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-stone-600">{t.heroLead}</p>
          </div>

          <div className="rounded-lg border border-stone-200 bg-white p-5 shadow-sm">
            <span
              className={`mb-4 block size-3 rounded-full ${
                apiStatus === "online" ? "bg-green-600" : "bg-amber-500"
              }`}
            />
            <strong>{apiStatusLabel(apiStatus, t)}</strong>
            <p className="mt-2 text-sm text-stone-600">{t.proxyTarget}: http://127.0.0.1:8080</p>
          </div>
        </section>

        <section id="events" className="border-t border-stone-200 py-12">
          <p className="text-sm font-bold uppercase text-moss">{t.eventsEyebrow}</p>
          <h2 className="mt-2 text-3xl font-black">{t.eventsTitle}</h2>
          <div className="mt-6 grid gap-5 md:grid-cols-2">
            {events.map((event) => (
              <article key={event.id} className="rounded-lg border border-stone-200 bg-white p-5">
                <div className="mb-4 h-36 rounded-md border border-stone-200 bg-green-50" />
                <h3 className="text-2xl font-black">{localized(event, "title", language)}</h3>
                <div className="mt-3 flex flex-wrap gap-3 text-sm text-stone-600">
                  <span>{event.date}</span>
                  <span>{event.time}</span>
                  <span>{event.venue}</span>
                  <span>{event.price}</span>
                </div>
                <p className="mt-4 leading-7 text-stone-600">
                  {localized(event, "summary", language)}
                </p>
              </article>
            ))}
          </div>
        </section>

        <section id="content" className="grid gap-7 border-t border-stone-200 py-12 md:grid-cols-[1fr_340px]">
          <div>
            <p className="text-sm font-bold uppercase text-moss">{t.contentEyebrow}</p>
            <h2 className="mt-2 text-3xl font-black">{t.contentTitle}</h2>
            <div className="mt-6 divide-y divide-stone-200">
              {contents.map((item) => (
                <div key={item.id} className="py-4">
                  <strong>{localized(item, "title", language)}</strong>
                  <p className="mt-1 text-stone-600">{localized(item, "summary", language)}</p>
                </div>
              ))}
            </div>
          </div>

          <aside id="connect" className="rounded-lg border border-stone-200 bg-white p-5">
            <p className="text-sm font-bold uppercase text-moss">{t.connectEyebrow}</p>
            <h2 className="mt-2 text-3xl font-black">{t.connectTitle}</h2>
            <div className="mt-5 grid gap-2">
              <LinkRow icon={<FileText size={18} />} label={t.linkArticles} />
              <LinkRow icon={<Camera size={18} />} label={t.linkPhotos} />
              <LinkRow icon={<Share2 size={18} />} label="SNS" />
              <LinkRow icon={<ShoppingBag size={18} />} label={t.linkShop} />
              <LinkRow icon={<Send size={18} />} label={t.linkSubmission} />
            </div>
          </aside>
        </section>

        <section id="join" className="border-t border-stone-200 py-12">
          <p className="text-sm font-bold uppercase text-moss">{t.joinEyebrow}</p>
          <h2 className="mt-2 text-3xl font-black">{t.joinTitle}</h2>
          <div className="mt-6 rounded-lg border border-stone-200 bg-white p-5 text-stone-600">
            {t.joinNoteBefore} <code>/api/join</code>
            {t.joinNoteAfter}
          </div>
        </section>
      </main>
    </div>
  );
}

function LanguageSwitch({
  language,
  onChange,
}: {
  language: Language;
  onChange: (language: Language) => void;
}) {
  return (
    <div className="inline-flex overflow-hidden rounded-lg border border-stone-300 bg-white" aria-label="语言选择">
      <button
        className={`px-4 py-2 text-sm font-bold ${language === "zh" ? "bg-moss text-white" : "text-stone-600"}`}
        type="button"
        aria-pressed={language === "zh"}
        onClick={() => onChange("zh")}
      >
        中文
      </button>
      <button
        className={`border-l border-stone-300 px-4 py-2 text-sm font-bold ${
          language === "ja" ? "bg-moss text-white" : "text-stone-600"
        }`}
        type="button"
        aria-pressed={language === "ja"}
        onClick={() => onChange("ja")}
      >
        日本語
      </button>
    </div>
  );
}

function apiStatusLabel(status: "checking" | "online" | "offline", t: Copy) {
  if (status === "online") return t.apiOnline;
  if (status === "offline") return t.apiOffline;
  return t.apiChecking;
}

function localized<T extends { title: string; summary: string }>(
  item: T & {
    titleI18n?: Partial<Record<Language, string>>;
    summaryI18n?: Partial<Record<Language, string>>;
  },
  field: "title" | "summary",
  language: Language,
) {
  const i18n = field === "title" ? item.titleI18n : item.summaryI18n;
  return i18n?.[language] || i18n?.zh || item[field];
}

function LinkRow({ icon, label }: { icon: ReactNode; label: string }) {
  return (
    <div className="flex items-center gap-3 border-b border-stone-200 py-3">
      {icon}
      <span>{label}</span>
    </div>
  );
}
