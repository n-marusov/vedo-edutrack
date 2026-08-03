/// <reference types="vite/client" />

// Build-time env fallbacks for the runtime config (ADR-DES.INFRA.dynamic-config-injection).
// The primary source is window.APP_CONFIG (/config.js, injected at deploy time);
// VITE_* values are the build-time fallback used by `make build` and local dev.
interface ImportMetaEnv {
  readonly VITE_APP_VERSION?: string;
  readonly VITE_APP_ENV?: string;
  readonly VITE_API_BASE_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
