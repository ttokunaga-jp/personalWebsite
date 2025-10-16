import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "node:path";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@shared": path.resolve(__dirname, "../../packages/shared/src")
    }
  },
  server: {
    port: 5174,
    proxy: {
      "/api": {
        target: "http://localhost:8100",
        changeOrigin: true
      }
    }
  },
  base: "/admin/"
});
