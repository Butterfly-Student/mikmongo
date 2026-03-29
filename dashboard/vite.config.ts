// dashboard/vite.config.ts
import path from "path"
import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import tailwindcss from "@tailwindcss/vite"
import { TanStackRouterVite } from "@tanstack/router-plugin/vite"

export default defineConfig({
  plugins: [
    // CRITICAL: TanStackRouterVite MUST be first — it scans src/routes/ before React transforms
    TanStackRouterVite({ target: "react", autoCodeSplitting: true }),
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/agent-portal": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/portal": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
})
