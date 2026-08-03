/**
 * Runtime application configuration — version, environment, URLs.
 *
 * Dynamic env injection (ADR-DES.INFRA.dynamic-config-injection):
 * the value is resolved at runtime in this order:
 *   1. window.APP_CONFIG  — set by /config.js (served/replaced at deploy:
 *      nginx envsubst template or Go-embed runtime generation);
 *   2. import.meta.env.VITE_* — build-time fallback (baked by `make build`);
 *   3. hardcoded dev defaults.
 */

export interface AppConfig {
  /** Build version (e.g. "1.2.3" or git describe output). */
  version: string;
  /** Environment label: development | staging | production. */
  env: string;
  /** Base URL of the EduTrack API (same-origin via the edge in production). */
  apiBaseUrl: string;
}

declare global {
  interface Window {
    APP_CONFIG?: Partial<AppConfig>;
  }
}

const runtimeConfig: Partial<AppConfig> =
  typeof window !== 'undefined' ? (window.APP_CONFIG ?? {}) : {};

export const appConfig: AppConfig = {
  version: runtimeConfig.version ?? import.meta.env.VITE_APP_VERSION ?? 'dev',
  env: runtimeConfig.env ?? import.meta.env.VITE_APP_ENV ?? 'development',
  apiBaseUrl: runtimeConfig.apiBaseUrl ?? import.meta.env.VITE_API_BASE_URL ?? '/api',
};
