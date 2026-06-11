import type { Language, LocalizedText } from "../types/api";

// 从实体里按当前语言取文案：优先当前语言的 i18n，回退中文，再回退非 i18n 的原字段。
// 后端一次性返回三语，前端只负责选值。
export function localized<T extends { [key: string]: unknown }>(item: T, field: string, language: Language) {
  const i18n = item[`${field}I18n`] as LocalizedText | undefined;
  const fallback = item[field];
  return i18n?.[language] || i18n?.zh || (typeof fallback === "string" ? fallback : "");
}

export function localizedMap(value: LocalizedText | undefined, language: Language) {
  return value?.[language] || value?.zh || "";
}
