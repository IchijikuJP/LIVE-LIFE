import type { Copy } from "../i18n";
import { languageLabels } from "../i18n";
import type { Language } from "../types/api";
import type { DesignVariant } from "../types/ui";

// 顶栏：品牌 + 导航 + 设计方案切换 + 语言切换。所有状态由父组件持有，这里只读 + 回调。
export function Header({
  copy,
  language,
  onLanguageChange,
  design,
  onDesignChange,
}: {
  copy: Copy;
  language: Language;
  onLanguageChange: (language: Language) => void;
  design: DesignVariant;
  onDesignChange: (design: DesignVariant) => void;
}) {
  return (
    <header className="topbar">
      <div className="topbar-actions">
        <a className="brand" href="#home" aria-label="LIVE LIFE">
          <span className="brand-mark" aria-hidden="true" />
          <nav aria-label="Primary navigation">
            <a href="#shows">{copy.navShows}</a>
            <a href="#cd-select">{copy.navCDSelect}</a>
            <a href="#archive">{copy.navArchive}</a>
            <a href="#connect">{copy.navConnect}</a>
          </nav>
        </a>

        <label className="design-switcher">
          <span>{copy.designLabel}</span>
          <select
            value={design}
            onChange={(event) =>
              onDesignChange(event.target.value as DesignVariant)
            }
          >
            <option value="v2">{copy.designV2}</option>
            <option value="v2-refined">{copy.designV2Refined}</option>
            <option value="v3">{copy.designV3}</option>
          </select>
        </label>

        <div className="language-switcher" aria-label="Language">
          {(["zh", "ja", "en"] as const).map((item) => (
            <button
              className={language === item ? "active" : ""}
              key={item}
              type="button"
              aria-pressed={language === item}
              onClick={() => onLanguageChange(item)}
            >
              {languageLabels[item]}
            </button>
          ))}
        </div>
      </div>
    </header>
  );
}
