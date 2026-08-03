import { appConfig } from '../../config';
import { useAuthStore } from '../../store/authStore';

/** ApiError — typed fetch failure with status code and server message. */
export class ApiError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

/**
 * Minimal typed fetch wrapper:
 *  - base URL from appConfig.apiBaseUrl
 *  - injects Authorization: Bearer <token> from the auth store
 *  - JSON request/response serialization
 *  - typed ApiError on non-2xx
 *  - auto-redirects to /login on 401
 */
async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = useAuthStore.getState().token;
  const headers = new Headers(init.headers);
  headers.set('Accept', 'application/json');
  if (init.body !== undefined && !(init.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  let res: Response;
  try {
    res = await fetch(`${appConfig.apiBaseUrl}${path}`, {
      ...init,
      headers,
    });
  } catch (err) {
    console.error('[api] network error', path, err);
    throw new ApiError(0, 'network_error', 'Network request failed');
  }

  if (res.status === 401) {
    console.warn('[api] 401 — redirecting to /login', path);
    useAuthStore.getState().logout();
    if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
      window.location.assign('/login');
    }
    throw new ApiError(401, 'unauthorized', 'Unauthorized');
  }

  if (!res.ok) {
    let code = 'request_failed';
    let message = `Request failed with status ${res.status}`;
    try {
      const body = (await res.json()) as { error?: string; message?: string };
      if (body.error) code = body.error;
      if (body.message) message = body.message;
    } catch {
      // non-JSON error body — keep defaults
    }
    console.error('[api] error', res.status, path, { code, message });
    throw new ApiError(res.status, code, message);
  }

  return (await res.json()) as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
};
