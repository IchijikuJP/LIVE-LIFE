const translations = {
  zh: {
    langAttr: "zh-Hans",
    navEvents: "活动 / 商店",
    navContent: "内容推荐",
    navJoin: "加入我们",
    heroEyebrow: "东京 Livehouse MVP",
    heroTitle: "演出情报、CD/唱片 Shop 和本地音乐现场入口。",
    heroLeadBefore: "这个本地预览页由 Go API 驱动。React 前端骨架已经放在",
    heroLeadAfter: "，下一步可以直接切换到正式前端开发。",
    healthChecking: "正在检查 API",
    healthOnline: "API 已连接",
    healthOffline: "API 未连接",
    eventsEyebrow: "下一场活动",
    eventsTitle: "活动 / 商店列表",
    contentEyebrow: "内容",
    contentTitle: "推荐内容",
    connectEyebrow: "连接",
    connectTitle: "关注与投稿",
    linkArticles: "文章",
    linkPhotos: "活动照片",
    linkVideos: "视频",
    linkSubmission: "投稿入口",
    joinEyebrow: "加入我们",
    joinTitle: "本地 API 表单测试",
    labelNickname: "昵称",
    labelEmail: "邮箱",
    labelRole: "身份",
    labelMessage: "留言",
    placeholderNickname: "LiveLife 朋友",
    placeholderEmail: "you@example.com",
    placeholderRole: "观众 / 音乐人 / 店铺 / 工作人员",
    placeholderMessage: "你想和 LiveLife 一起做什么？",
    submitButton: "提交",
    submitting: "提交中...",
    submitFallback: "已收到你的本地测试报名。",
    submitFailed: "提交失败",
    requestFailed: "请求失败",
  },
  ja: {
    langAttr: "ja",
    navEvents: "イベント / ショップ",
    navContent: "おすすめ",
    navJoin: "参加する",
    heroEyebrow: "東京ライブハウス MVP",
    heroTitle: "ライブ情報、CD/レコードショップ、ローカルシーンへの入口。",
    heroLeadBefore: "このローカルプレビューは Go API で動いています。React フロントエンドの骨組みは",
    heroLeadAfter: "に用意してあるので、次の段階で正式なフロント開発へ移れます。",
    healthChecking: "API を確認中",
    healthOnline: "API 接続済み",
    healthOffline: "API 未接続",
    eventsEyebrow: "次のイベント",
    eventsTitle: "イベント / ショップ一覧",
    contentEyebrow: "コンテンツ",
    contentTitle: "おすすめ",
    connectEyebrow: "つながる",
    connectTitle: "フォローと投稿",
    linkArticles: "記事",
    linkPhotos: "イベント写真",
    linkVideos: "動画",
    linkSubmission: "投稿入口",
    joinEyebrow: "参加する",
    joinTitle: "ローカル API フォームテスト",
    labelNickname: "ニックネーム",
    labelEmail: "メール",
    labelRole: "立場",
    labelMessage: "メッセージ",
    placeholderNickname: "LiveLife の友人",
    placeholderEmail: "you@example.com",
    placeholderRole: "観客 / アーティスト / ショップ / スタッフ",
    placeholderMessage: "LiveLife と一緒に何をしたいですか？",
    submitButton: "送信",
    submitting: "送信中...",
    submitFallback: "ローカルテストの参加リクエストを受け取りました。",
    submitFailed: "送信に失敗しました",
    requestFailed: "リクエストに失敗しました",
  },
};

const healthStatus = document.querySelector("#healthStatus");
const healthMeta = document.querySelector("#healthMeta");
const statusDot = document.querySelector(".status-dot");
const eventsGrid = document.querySelector("#eventsGrid");
const contentList = document.querySelector("#contentList");
const joinForm = document.querySelector("#joinForm");
const formResult = document.querySelector("#formResult");
const languageButtons = Array.from(document.querySelectorAll(".language-button"));

let currentLanguage = "zh";
let cachedEvents = [];
let cachedContents = [];
let lastHealth = null;

function t(key) {
  return translations[currentLanguage][key] || translations.zh[key] || key;
}

function localizedText(item, field) {
  const i18n = item[`${field}I18n`];
  return i18n?.[currentLanguage] || i18n?.zh || item[field] || "";
}

async function getJSON(url) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`${t("requestFailed")}: ${response.status}`);
  }
  return response.json();
}

function applyLanguage(language) {
  currentLanguage = translations[language] ? language : "zh";
  document.documentElement.lang = t("langAttr");

  document.querySelectorAll("[data-i18n]").forEach((node) => {
    node.textContent = t(node.dataset.i18n);
  });

  document.querySelectorAll("[data-i18n-placeholder]").forEach((node) => {
    node.setAttribute("placeholder", t(node.dataset.i18nPlaceholder));
  });

  languageButtons.forEach((button) => {
    const active = button.dataset.lang === currentLanguage;
    button.classList.toggle("active", active);
    button.setAttribute("aria-pressed", String(active));
  });

  if (!lastHealth) {
    healthStatus.textContent = t("healthChecking");
  } else if (lastHealth.ok) {
    healthStatus.textContent = t("healthOnline");
  } else {
    healthStatus.textContent = t("healthOffline");
  }

  renderEvents(cachedEvents);
  renderContents(cachedContents);
}

function renderEvent(event) {
  const tags = event.tags.map((tag) => `<span class="tag">${tag}</span>`).join("");
  return `
    <article class="card">
      <div class="poster" aria-hidden="true"></div>
      <h3>${localizedText(event, "title")}</h3>
      <div class="meta">
        <span>${event.date}</span>
        <span>${event.time}</span>
        <span>${event.venue}</span>
        <span>${event.price}</span>
      </div>
      <p class="summary">${localizedText(event, "summary")}</p>
      <div class="tags">${tags}</div>
    </article>
  `;
}

function renderContent(item) {
  return `
    <div class="content-item">
      <strong>${localizedText(item, "title")}</strong>
      <span>${localizedText(item, "summary")}</span>
    </div>
  `;
}

function renderEvents(events) {
  eventsGrid.innerHTML = events.map(renderEvent).join("");
}

function renderContents(contents) {
  contentList.innerHTML = contents.map(renderContent).join("");
}

async function boot() {
  applyLanguage(currentLanguage);

  try {
    const health = await getJSON("/api/health");
    lastHealth = { ok: true, value: health };
    healthStatus.textContent = t("healthOnline");
    healthMeta.textContent = `${health.service} / ${health.status}`;
    statusDot.classList.add("ok");
  } catch (error) {
    lastHealth = { ok: false, value: error };
    healthStatus.textContent = t("healthOffline");
    healthMeta.textContent = error.message;
  }

  try {
    const [{ events }, { contents }] = await Promise.all([
      getJSON("/api/events"),
      getJSON("/api/contents"),
    ]);
    cachedEvents = events;
    cachedContents = contents;
    renderEvents(cachedEvents);
    renderContents(cachedContents);
  } catch (error) {
    eventsGrid.innerHTML = `<p class="summary">${error.message}</p>`;
  }
}

languageButtons.forEach((button) => {
  button.addEventListener("click", () => applyLanguage(button.dataset.lang));
});

joinForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  formResult.textContent = t("submitting");

  const data = Object.fromEntries(new FormData(joinForm).entries());

  try {
    const response = await fetch("/api/join", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    const result = await response.json();
    if (!response.ok) {
      throw new Error(result.error || t("submitFailed"));
    }
    formResult.textContent = result.message || t("submitFallback");
    joinForm.reset();
  } catch (error) {
    formResult.textContent = error.message;
  }
});

boot();
