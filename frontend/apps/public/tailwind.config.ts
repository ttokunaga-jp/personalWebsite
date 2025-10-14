import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./index.html", "./src/**/*.{ts,tsx}", "../../packages/shared/src/**/*.{ts,tsx}"],
  theme: {
    extend: {}
  },
  plugins: []
};

export default config;
