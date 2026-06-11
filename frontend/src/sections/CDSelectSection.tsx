import { useMemo, useState } from "react";

import { CatalogCard } from "../components/CatalogCard";
import type { Copy } from "../i18n";
import type { CatalogItem, Language } from "../types/api";

type FormatFilter = "all" | "cd" | "vinyl";

// CD 严选：内部分 CD / 黑胶，格式筛选是本板块自己的本地 UI 状态。
export function CDSelectSection({
  copy,
  language,
  catalog,
}: {
  copy: Copy;
  language: Language;
  catalog: CatalogItem[];
}) {
  const [formatFilter, setFormatFilter] = useState<FormatFilter>("all");
  const filteredCatalog = useMemo(
    () => (formatFilter === "all" ? catalog : catalog.filter((item) => item.format === formatFilter)),
    [catalog, formatFilter],
  );

  return (
    <section id="cd-select" className="section">
      <div className="section-heading split-heading">
        <div>
          <p className="eyebrow">{copy.cdEyebrow}</p>
          <h2>{copy.cdTitle}</h2>
        </div>
        <p>{copy.cdCopy}</p>
      </div>

      <div className="format-tabs" aria-label="CD select formats">
        {([
          ["all", copy.formatAll],
          ["cd", copy.formatCD],
          ["vinyl", copy.formatVinyl],
        ] as const).map(([value, label]) => (
          <button
            key={value}
            className={formatFilter === value ? "active" : ""}
            type="button"
            onClick={() => setFormatFilter(value)}
          >
            {label}
          </button>
        ))}
      </div>

      <div className="catalog-grid">
        {filteredCatalog.map((item) => (
          <CatalogCard key={item.id} item={item} language={language} copy={copy} />
        ))}
      </div>
    </section>
  );
}
