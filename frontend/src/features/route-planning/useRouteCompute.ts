import { useCallback, useState } from 'react';
import type { RouterTopic } from '../../shared/api/types';
import { computeRoute } from './api';

/** Horizon view of the computed route. */
export type Horizon = 'far' | 'mid' | 'near';

export interface RouteComputeState {
  route: RouterTopic[];
  horizon: Horizon;
  loading: boolean;
  error: string | null;
  recomputeNeeded: boolean;
}

export interface RouteBuilderState extends RouteComputeState {
  setHorizon: (horizon: Horizon) => void;
  compute: (goalTopicId: string) => Promise<void>;
}

const midHorizonCount = 5;

/** useRouteCompute — triggers POST /routes/compute and manages horizon views. */
export function useRouteCompute(learnerId: string): RouteBuilderState {
  const [state, setState] = useState<RouteComputeState>({
    route: [],
    horizon: 'far',
    loading: false,
    error: null,
    recomputeNeeded: false,
  });

  const compute = useCallback(
    async (goalTopicId: string) => {
      if (!goalTopicId.trim()) {
        setState((prev) => ({ ...prev, error: 'Goal topic is required' }));
        return;
      }
      setState((prev) => ({ ...prev, loading: true, error: null }));
      try {
        const route = await computeRoute(learnerId, goalTopicId.trim());
        setState((prev) => ({
          ...prev,
          route,
          loading: false,
          error: null,
          recomputeNeeded: false,
        }));
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        setState((prev) => ({ ...prev, loading: false, error: message }));
      }
    },
    [learnerId],
  );

  const setHorizon = useCallback((horizon: Horizon) => {
    setState((prev) => ({ ...prev, horizon }));
  }, []);

  return { ...state, setHorizon, compute };
}

/** Apply the current horizon slice to the full route. */
export function horizonSlice(route: RouterTopic[], horizon: Horizon): RouterTopic[] {
  if (horizon === 'near') {
    return route.slice(0, 1);
  }
  if (horizon === 'mid') {
    return route.slice(0, midHorizonCount);
  }
  return route;
}
