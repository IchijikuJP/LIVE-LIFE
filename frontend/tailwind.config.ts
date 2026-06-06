import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#111816",
        moss: "#23613d",
        paper: "#f6f7f4",
      },
    },
  },
  plugins: [],
} satisfies Config;
