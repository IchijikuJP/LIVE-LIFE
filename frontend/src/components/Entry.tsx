import type { ReactNode } from "react";

// 首页四个板块入口磁贴（演出 / CD / 档案 / 联系）。
export function Entry({ href, index, label, icon }: { href: string; index: string; label: string; icon: ReactNode }) {
  return (
    <a className="entry-tile" href={href}>
      <span>
        <strong>{index}</strong>
        {icon}
      </span>
      <strong>{label}</strong>
    </a>
  );
}
