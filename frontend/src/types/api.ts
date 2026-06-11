// LIVE LIFE 前后端 API 契约类型。
// 这里是前端对后端 /api/* 返回结构的单一事实来源，对应 Go 端 internal/domain 的实体。
// 后端改字段时只改这里，避免类型散落在各组件里各自维护、悄悄 drift。

// API 一次性返回三语言文案，前端按当前语言取值。键对应后端 LocalizedText 的 zh/ja/en。
export type Language = "zh" | "ja" | "en";
export type LocalizedText = Partial<Record<Language, string>>;

// 对应 domain.Event。owned=true 表示 LIVE LIFE 自主演出，前端固定置顶。
export type EventItem = {
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

// 对应 domain.CatalogItem。format 把 CD 严选拆成 cd / vinyl，购买走 purchaseUrl 跳外部 Shop。
export type CatalogItem = {
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

// 对应 domain.ContentItem，用于 Archive 档案馆内容。
export type ContentItem = {
  id: string;
  type: string;
  title: string;
  titleI18n?: LocalizedText;
  summary: string;
  summaryI18n?: LocalizedText;
};

// 对应 GET /api/health 返回。
export type HealthPayload = {
  status: string;
  brand: string;
  service: string;
};
