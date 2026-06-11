import type { Copy } from "../i18n";
import { localized } from "../lib/localized";
import type { ContentItem, Language } from "../types/api";

// 档案馆：历史海报、公开资料备注等内容卡片。
export function ArchiveSection({
  copy,
  language,
  contents,
}: {
  copy: Copy;
  language: Language;
  contents: ContentItem[];
}) {
  return (
    <section id="archive" className="section">
      <div className="section-heading split-heading">
        <div>
          <p className="eyebrow">{copy.archiveEyebrow}</p>
          <h2>{copy.archiveTitle}</h2>
        </div>
        <p>{copy.archiveCopy}</p>
      </div>

      <div className="grid compact-grid">
        {contents.map((item) => (
          <article className="content-card" key={item.id}>
            <span className="pill">{item.type}</span>
            <h3>{localized(item, "title", language)}</h3>
            <p>{localized(item, "summary", language)}</p>
          </article>
        ))}
      </div>
    </section>
  );
}
