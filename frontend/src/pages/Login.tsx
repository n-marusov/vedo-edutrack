import { useState } from 'react';
import type { FormEvent } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { ApiError, api } from '../shared/api';
import type { TokenResponse } from '../shared/api/types';
import { Button } from '../shared/components/Button';
import { Input } from '../shared/components/Input';
import { useAuthStore } from '../store/authStore';

const roleOptions = ['learner', 'parent', 'teacher', 'methodologist', 'admin'];

/**
 * LoginPage — dev/demo auth flow: user_id + role → JWT.
 * Real credential login lands post-MVP (Keycloak).
 */
export function LoginPage() {
  const [userId, setUserId] = useState('demo');
  const [role, setRole] = useState('learner');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const login = useAuthStore((s) => s.login);
  const navigate = useNavigate();
  const location = useLocation();

  const from = (location.state as { from?: string } | null)?.from ?? '/dashboard';

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    console.info('[login] attempt', { userId });
    try {
      const res = await api.post<TokenResponse>('/auth/token', {
        user_id: userId,
        roles: [role],
      });
      login(res.access_token, { userId, roles: [role] });
      navigate(from, { replace: true });
    } catch (err) {
      console.warn('[login] failed', err);
      setError(err instanceof ApiError ? err.message : 'Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h1 className="text-xl font-semibold text-gray-900">Sign in</h1>
      <p className="text-sm text-gray-500">Demo login: choose a user ID and role.</p>

      <Input
        label="User ID"
        value={userId}
        onChange={(e) => setUserId(e.target.value)}
        placeholder="demo"
        autoComplete="username"
      />

      <div>
        <label htmlFor="role" className="mb-1 block text-sm font-medium text-gray-700">
          Role
        </label>
        <select
          id="role"
          value={role}
          onChange={(e) => setRole(e.target.value)}
          className="block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          {roleOptions.map((r) => (
            <option key={r} value={r}>
              {r}
            </option>
          ))}
        </select>
      </div>

      {error && (
        <p className="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700" role="alert">
          {error}
        </p>
      )}

      <Button type="submit" loading={loading} className="w-full">
        Sign in
      </Button>
    </form>
  );
}
