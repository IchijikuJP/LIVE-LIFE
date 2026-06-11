import type { Language } from "../types/api";

import { en } from "./en";
import { ja } from "./ja";
import type { Copy } from "./types";
import { zh } from "./zh";

export type { Copy };

// 默认中文。三语文案在这里汇总，组件按当前 language 取 copy[language]。
export const copy: Record<Language, Copy> = { zh, ja, en };

export const languageLabels: Record<Language, string> = {
  zh: "中文",
  ja: "日本語",
  en: "ENGLISH",
};
