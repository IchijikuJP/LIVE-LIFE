import type { CatalogItem, ContentItem, EventItem, HealthPayload } from "../types/api";

// 后端列表接口都包了一层 brand + 分组字段（ownedEvents/cd/vinyl 等），
// 前端首页只取完整列表，分组逻辑在前端用 owned/format 自己过滤。
type EventsResponse = { events?: EventItem[] };
type CDItemsResponse = { items?: CatalogItem[] };
type ContentsResponse = { contents?: ContentItem[] };

async function getJSON<T>(path: string): Promise<T> {
  const res = await fetch(path);
  if (!res.ok) {
    throw new Error(`API ${path} failed: ${res.status}`);
  }
  return (await res.json()) as T;
}

export type LiveLifeData = {
  health: HealthPayload;
  events: EventItem[];
  catalog: CatalogItem[];
  contents: ContentItem[];
};

// 首页启动时一次性并行拉取四个只读接口。任一失败则整体 reject，由调用方标记 offline。
export async function fetchLiveLifeData(): Promise<LiveLifeData> {
  const [health, events, cd, contents] = await Promise.all([
    getJSON<HealthPayload>("/api/health"),
    getJSON<EventsResponse>("/api/events"),
    getJSON<CDItemsResponse>("/api/cd-items"),
    getJSON<ContentsResponse>("/api/contents"),
  ]);
  return {
    health,
    events: events.events ?? [],
    catalog: cd.items ?? [],
    contents: contents.contents ?? [],
  };
}
