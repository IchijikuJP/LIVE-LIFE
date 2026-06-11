// V2 设计方案的背景纹理：滚动的经典摇滚乐队名 + 抽象音轨场。
// 纯装饰（aria-hidden），不含业务数据。显隐由 body[data-design] 的 CSS 控制。
const v2BandNames = [
  "THE BEATLES",
  "THE ROLLING STONES",
  "THE KINKS",
  "THE WHO",
  "THE YARDBIRDS",
  "CREAM",
  "LED ZEPPELIN",
  "PINK FLOYD",
  "BLACK SABBATH",
  "DEEP PURPLE",
  "QUEEN",
  "DAVID BOWIE",
  "T. REX",
  "ROXY MUSIC",
  "SPARKS",
  "THE VELVET UNDERGROUND",
  "THE STOOGES",
  "MC5",
  "NEW YORK DOLLS",
  "PATTI SMITH",
  "TELEVISION",
  "RAMONES",
  "BLONDIE",
  "TALKING HEADS",
  "THE CLASH",
  "SEX PISTOLS",
  "BUZZCOCKS",
  "WIRE",
  "XTC",
  "THE JAM",
  "PUBLIC IMAGE LTD",
  "JOY DIVISION",
  "NEW ORDER",
  "THE FALL",
  "GANG OF FOUR",
  "SIOUXSIE AND THE BANSHEES",
  "BAUHAUS",
  "THE CURE",
  "ECHO AND THE BUNNYMEN",
  "THE PSYCHEDELIC FURS",
  "THE SMITHS",
  "THE JESUS AND MARY CHAIN",
  "MY BLOODY VALENTINE",
  "SLOWDIVE",
  "RIDE",
  "LUSH",
  "COCTEAU TWINS",
  "SPACEMEN 3",
  "THE STONE ROSES",
  "HAPPY MONDAYS",
  "PRIMAL SCREAM",
  "MASSIVE ATTACK",
  "PORTISHEAD",
  "TRICKY",
  "RADIOHEAD",
  "OASIS",
  "BLUR",
  "SUEDE",
  "PULP",
  "THE VERVE",
  "MOGWAI",
  "PJ HARVEY",
  "NIRVANA",
  "PEARL JAM",
  "SONIC YOUTH",
  "PIXIES",
  "DINOSAUR JR.",
  "R.E.M.",
  "THE REPLACEMENTS",
  "HUSKER DU",
  "GUIDED BY VOICES",
  "PAVEMENT",
  "YO LA TENGO",
  "THE FLAMING LIPS",
  "WILCO",
  "THE WHITE STRIPES",
  "THE STROKES",
  "INTERPOL",
  "YEAH YEAH YEAHS",
  "LCD SOUNDSYSTEM",
  "ARCADE FIRE",
  "TV ON THE RADIO",
  "ARCTIC MONKEYS",
  "FRANZ FERDINAND",
  "BLOC PARTY",
  "THE LIBERTINES",
  "FOALS",
  "IDLES",
  "FONTAINES D.C.",
  "BLACK MIDI",
  "SQUID",
  "CAN",
  "NEU!",
  "KRAFTWERK",
  "FAUST",
  "TALK TALK",
  "WIRE",
  "SWANS",
  "THE RESIDENTS",
  "SUICIDE",
  "DEVO",
  "MISSION OF BURMA",
  "MINUTEMEN",
  "FUGAZI",
  "BIG BLACK",
  "SLINT",
  "TORTOISE",
  "STEREOLAB",
  "BROADCAST",
  "THE SEA AND CAKE",
];

export function V2Texture() {
  return (
    <>
      <V2BandScroll />
      <div className="band-texture" aria-hidden="true">
        THE SMITHS OASIS BLUR RADIOHEAD PULP SUEDE JOY DIVISION NEW ORDER THE CURE STONE ROSES MY BLOODY VALENTINE SLOWDIVE PRIMAL SCREAM MASSIVE ATTACK PORTISHEAD SONIC YOUTH PIXIES
      </div>
      <div className="track-field" aria-hidden="true">
        <span className="lane lane-a" />
        <span className="lane lane-b" />
        <span className="lane lane-c" />
        <span className="sample sample-a" />
        <span className="sample sample-b" />
        <span className="wave wave-a" />
        <span className="wave wave-b" />
      </div>
    </>
  );
}

function V2BandScroll() {
  const columns = [
    v2BandNames,
    [...v2BandNames].slice(28).concat(v2BandNames.slice(0, 28)),
    [...v2BandNames].slice(58).concat(v2BandNames.slice(0, 58)),
  ];

  return (
    <div className="v2-band-scroll" aria-hidden="true">
      {columns.map((names, index) => (
        <div className={`v2-band-column column-${index + 1}`} key={index}>
          {[...names, ...names].map((name, itemIndex) => (
            <span key={`${name}-${itemIndex}`}>{name}</span>
          ))}
        </div>
      ))}
    </div>
  );
}
