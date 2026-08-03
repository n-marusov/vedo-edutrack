import { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';
import { RouterProvider } from 'react-router';
import { AuthProvider } from './features/identity-access';
import { router } from './routes';
import { ErrorFallback } from './shared/components';

/** Root error boundary — catches render errors and shows a retry UI. */
class RootErrorBoundary extends Component<{ children: ReactNode }, { error: Error | null }> {
  state = { error: null as Error | null };

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('[app] root error boundary', error, info.componentStack);
  }

  render() {
    if (this.state.error) {
      return (
        <ErrorFallback
          error={this.state.error}
          resetErrorBoundary={() => this.setState({ error: null })}
        />
      );
    }
    return this.props.children;
  }
}

/**
 * VEDO EduTrack root application component.
 *
 * Wraps the router in the auth provider (token validation on mount) and a
 * catch-all error boundary. Lazy pages render inside Suspense (see routes.tsx).
 */
export function App() {
  return (
    <RootErrorBoundary>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </RootErrorBoundary>
  );
}
