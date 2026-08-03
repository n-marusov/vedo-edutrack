// Runtime application configuration (dynamic env injection).
//
// ADR-DES.INFRA.dynamic-config-injection: one build → many environments.
// The app reads this object via src/config.ts (window.APP_CONFIG), preferring
// runtime values over build-time import.meta.env.VITE_* fallbacks.
//
// This file is the DEV default (served as-is by Vite from public/).
// In deployed variants it is REPLACED at runtime:
//   - nginx (Dockerfile.nginx): /etc/nginx/templates/config.js.template
//     is env-substituted at container start (APP_VERSION, APP_ENV, API_BASE_URL)
//   - Go embed (Dockerfile.embed): the embed server generates /config.js
//     from env vars at runtime
//
// `version` is intentionally omitted here — in dev it falls back to the
// build-time/default ("dev") so local builds show their real VITE_APP_VERSION.
window.APP_CONFIG = {
  env: 'development',
  apiBaseUrl: '/api',
};
