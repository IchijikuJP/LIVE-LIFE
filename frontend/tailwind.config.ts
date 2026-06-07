import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        paper: "#f4efdf",
        ink: "#111111",
        muted: "#665f53",
        yellow: "#ffd000",
        blue: "#2457ff",
        red: "#e5002a",
        green: "#00c26f",
      },
    },
  },
  plugins: [],
} satisfies Config;
