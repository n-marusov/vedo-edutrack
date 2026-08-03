import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

// Vitest configuration — VEDO EduTrack frontend (M0.2, T14).
// Component/unit tests run in jsdom via Vitest + React Testing Library
// (ADR-IMPL.PROCESS.development-tooling §6).

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/__tests__/**/*.test.{ts,tsx}'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['src/**/*.{ts,tsx}'],
      // M0.2: advisory — thresholds are enforced from M1
      // (REQ-NFR-process.dev.test-coverage, core >= 90%).
    },
  },
});
