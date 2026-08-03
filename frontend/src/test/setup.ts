import '@testing-library/jest-dom/vitest';
import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';

// RTL cleanup between tests (Vitest + RTL, ADR-IMPL.PROCESS.development-tooling §6).
afterEach(() => {
  cleanup();
});

// jsdom does not implement matchMedia — components that use responsive
// breakpoints (M0.3+) rely on it. Minimal no-op mock.
if (typeof window !== 'undefined' && !window.matchMedia) {
  window.matchMedia = (query: string): MediaQueryList =>
    ({
      matches: false,
      media: query,
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    }) as unknown as MediaQueryList;
}
