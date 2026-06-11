import { FormEvent, useEffect, useMemo, useState } from "react";

import { copy } from "./i18n";
import type { Copy } from "./i18n";
import { useLiveLifeData } from "./hooks/useLiveLifeData";
import type { Language } from "./types/api";
import type { DesignVariant } from "./types/ui";
import { Header } from "./sections/Header";
import { HeroSection } from "./sections/HeroSection";
import { EntryGrid } from "./sections/EntryGrid";
import { ShowsSection } from "./sections/ShowsSection";
import { CDSelectSection } from "./sections/CDSelectSection";
import { ArchiveSection } from "./sections/ArchiveSection";
import { ConnectSection } from "./sections/ConnectSection";

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

function readInitialDesign(): DesignVariant {
  const params = new URLSearchParams(window.location.search);
  const fromUrl = params.get("design");
  if (fromUrl === "v2" || fromUrl === "v2-refined" || fromUrl === "v3") {
    return fromUrl;
  }
  const stored = window.localStorage.getItem("liveLifeDesignVariant");
  return stored === "v3" || stored === "v2-refined" ? stored : "v2";
}

// App 现在只负责：全局状态、派生数据、副作用、Connect 提交编排，以及把各 section 组合起来。
// 具体展示都在 sections/ 与 components/ 里，设计纹理在 design/，数据加载在 hooks/。
export function App() {
  const [language, setLanguage] = useState<Language>("zh");
  const [design, setDesign] = useState<DesignVariant>(readInitialDesign);
  const [formStatus, setFormStatus] = useState("");
  const { events, catalog, contents, health, apiStatus } = useLiveLifeData();
  const t = copy[language];

  const ownedEvents = useMemo(() => events.filter((event) => event.owned), [events]);
  const recommendedEvents = useMemo(() => events.filter((event) => !event.owned), [events]);

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
      <Header
        copy={t}
        language={language}
        onLanguageChange={setLanguage}
        design={design}
        onDesignChange={setDesign}
      />

      <main>
        <HeroSection
          copy={t}
          language={language}
          apiStatus={apiStatus}
          health={health}
          ownedEvents={ownedEvents}
          catalog={catalog}
        />
        <EntryGrid copy={t} />
        <ShowsSection
          copy={t}
          language={language}
          ownedEvents={ownedEvents}
          recommendedEvents={recommendedEvents}
        />
        <CDSelectSection copy={t} language={language} catalog={catalog} />
        <ArchiveSection copy={t} language={language} contents={contents} />
        <ConnectSection copy={t} formStatus={formStatus} onSubmit={submitConnect} />
      </main>
    </div>
  );
}
