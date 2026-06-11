import { useEffect, useState } from "react";

import { fetchLiveLifeData } from "../api/client";
import type { CatalogItem, ContentItem, EventItem, HealthPayload } from "../types/api";

export type ApiStatus = "checking" | "online" | "offline";

export type LiveLifeDataState = {
  events: EventItem[];
  catalog: CatalogItem[];
  contents: ContentItem[];
  health: HealthPayload | null;
  apiStatus: ApiStatus;
};

// 启动后一次性加载首页只读数据。组件只读结果，不关心 fetch / 端点 / 错误处理细节。
// 后端换 PostgreSQL 或后台管理系统时，只要 API 契约不变，这里和组件都不用改。
export function useLiveLifeData(): LiveLifeDataState {
  const [events, setEvents] = useState<EventItem[]>([]);
  const [catalog, setCatalog] = useState<CatalogItem[]>([]);
  const [contents, setContents] = useState<ContentItem[]>([]);
  const [health, setHealth] = useState<HealthPayload | null>(null);
  const [apiStatus, setApiStatus] = useState<ApiStatus>("checking");

  useEffect(() => {
    let alive = true;
    fetchLiveLifeData()
      .then((data) => {
        if (!alive) return;
        setHealth(data.health);
        setEvents(data.events);
        setCatalog(data.catalog);
        setContents(data.contents);
        setApiStatus("online");
      })
      .catch(() => {
        if (alive) setApiStatus("offline");
      });
    return () => {
      alive = false;
    };
  }, []);

  return { events, catalog, contents, health, apiStatus };
}
