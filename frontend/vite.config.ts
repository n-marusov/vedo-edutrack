import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// Vite configuration for VEDO EduTrack React SPA.
// See ADR-DES.STACK.framework-vs-vs (React + Vite).

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": import.meta.dirname + "/src",
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        // Containerized dev (docker-compose): the frontend container reaches
        // the backend via the service name; local dev keeps localhost.
        target: process.env.VITE_PROXY_TARGET ?? "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    sourcemap: true,
  },
});
