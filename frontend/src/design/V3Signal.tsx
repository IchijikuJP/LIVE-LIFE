// V3「乐队信号」设计方案的背景：密集乐队名云 + LIVE LIFE 字母轮播 + 抽象音轨。
// 纯装饰（aria-hidden），不含业务数据。显隐由 body[data-design] 的 CSS 控制。
const classicBandCloud =
  "THE BEATLES THE ROLLING STONES LED ZEPPELIN PINK FLOYD QUEEN THE WHO DAVID BOWIE THE KINKS T. REX ROXY MUSIC THE CLASH SEX PISTOLS THE JAM BUZZCOCKS JOY DIVISION NEW ORDER THE SMITHS THE CURE SIOUXSIE AND THE BANSHEES ECHO AND THE BUNNYMEN THE STONE ROSES MY BLOODY VALENTINE SLOWDIVE PRIMAL SCREAM RADIOHEAD OASIS BLUR SUEDE PULP THE VERVE ARCTIC MONKEYS NIRVANA PIXIES SONIC YOUTH R.E.M. THE STOOGES THE VELVET UNDERGROUND TALKING HEADS";

export function V3Signal() {
  return (
    <>
      <div className="v3-brand-signal" aria-hidden="true">
        <div className="v3-logo-mark" />
        <div className="v3-band-cloud">{classicBandCloud}</div>
        <div className="v3-marquee-row row-one">LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE LIVE LIFE</div>
        <div className="v3-marquee-row row-two">L I V E L I F E L I V E L I F E L I V E L I F E L I V E L I F E</div>
      </div>
      <div className="v3-track-field" aria-hidden="true">
        <div className="song-timeline">
          <span>THE LIBERTINES / MUSIC WHEN THE LIGHTS GO OUT</span>
          <span>126 BPM</span>
          <span>INDIE ROCK SESSION</span>
        </div>
        <span className="section-marker marker-a" />
        <span className="section-marker marker-b" />
        <span className="section-marker marker-c" />
        <span className="track-wave mix-wave-a" />
        <span className="track-row row-a" />
        <span className="track-row row-b" />
        <span className="track-row row-c" />
        <span className="track-row row-d" />
        <span className="track-block block-a" />
        <span className="track-block block-b" />
        <span className="track-block block-c" />
        <span className="beat-grid beat-a" />
        <span className="beat-grid beat-b" />
        <span className="track-clip clip-a" />
        <span className="track-clip clip-b" />
        <span className="track-clip clip-c" />
        <span className="track-wave mix-wave-b" />
        <span className="track-meter meter-a" />
        <span className="track-meter meter-b" />
      </div>
    </>
  );
}
