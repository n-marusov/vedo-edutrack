import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { App } from './App';
import { initTelemetry } from './shared/telemetry';
import './styles/index.css';

initTelemetry();

const root = document.getElementById('root');
if (!root) {
  throw new Error('Root element not found — check index.html');
}

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
