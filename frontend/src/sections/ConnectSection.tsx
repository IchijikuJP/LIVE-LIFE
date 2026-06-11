import type { FormEvent } from "react";

import type { Copy } from "../i18n";

// 统一联系入口：票务 / 购买售后 / 合作 / 投稿。
// 这里只负责表单展示，提交逻辑（mailto + POST）暂由父组件传入，后续抽成 useConnectForm。
export function ConnectSection({
  copy,
  formStatus,
  onSubmit,
}: {
  copy: Copy;
  formStatus: string;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
}) {
  return (
    <section id="connect" className="section connect-section">
      <div className="connect-copy">
        <p className="eyebrow">{copy.connectEyebrow}</p>
        <h2>{copy.connectTitle}</h2>
        <p>{copy.connectCopy}</p>
      </div>

      <form className="connect-form" onSubmit={onSubmit}>
        <label>
          <span>{copy.labelNickname}</span>
          <input name="nickname" placeholder={copy.placeholderNickname} required />
        </label>
        <label>
          <span>{copy.labelEmail}</span>
          <input name="email" type="email" placeholder={copy.placeholderEmail} required />
        </label>
        <label>
          <span>{copy.labelTopic}</span>
          <select name="topic" defaultValue="ticket">
            <option value="ticket">{copy.topicTicket}</option>
            <option value="cd-select">{copy.topicCDSelect}</option>
            <option value="support">{copy.topicSupport}</option>
            <option value="collab">{copy.topicCollab}</option>
          </select>
        </label>
        <label className="wide">
          <span>{copy.labelMessage}</span>
          <textarea name="message" rows={5} placeholder={copy.placeholderMessage} />
        </label>
        <button type="submit">{copy.submitButton}</button>
        <p className="form-result" role="status">{formStatus}</p>
      </form>
    </section>
  );
}
