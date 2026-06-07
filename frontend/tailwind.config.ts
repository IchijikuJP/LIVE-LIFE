import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        night: "#111111",
        cream: "#f5f1ea",
        amber: "#f5c84b",
        pink: "#e04f9a",
        green: "#7fd39a",
      },
    },
  },
  plugins: [],
} satisfies Config;
