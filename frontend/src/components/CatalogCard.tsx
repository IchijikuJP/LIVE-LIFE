import type { Copy } from "../i18n";
import { localized, localizedMap } from "../lib/localized";
import type { CatalogItem, Language } from "../types/api";

// CD 严选单品卡片。购买不走站内，点 buy-button 跳到 purchaseUrl（BASE 等外部 Shop）。
export function CatalogCard({ item, language, copy }: { item: CatalogItem; language: Language; copy: Copy }) {
  return (
    <article className="catalog-card">
      {item.imageUrl ? (
        <img src={item.imageUrl} alt={localized(item, "title", language)} />
      ) : (
        <div className="catalog-art" aria-hidden="true" />
      )}
      <div className="catalog-body">
        <span className="pill">{item.format === "vinyl" ? copy.formatVinyl : copy.formatCD}</span>
        <h3>{localized(item, "title", language)}</h3>
        <p className="artist">{item.artist}</p>
        <p>{localized(item, "summary", language)}</p>
        <div className="track-list">{item.tracks?.map((track) => <span key={track}>{track}</span>)}</div>
        <div className="catalog-actions">
          <span>{item.price || "TBD"}</span>
          <a className="button buy-button" href={item.purchaseUrl} target="_blank" rel="noreferrer">
            {localizedMap(item.purchaseText, language) || "BUY HERE"}
          </a>
        </div>
        <small>{copy.externalShopNote}</small>
      </div>
    </article>
  );
}
