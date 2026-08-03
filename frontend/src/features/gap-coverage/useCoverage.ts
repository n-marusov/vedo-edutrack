import { useEffect, useState } from 'react';
import type { CoverageResponse, GapDiagnosisResponse } from '../../shared/api/types';
import { diagnoseGaps, getFgosCoverage } from '../execution-progress/api';

export interface GapCoverageState {
  coverage: CoverageResponse | null;
  gaps: GapDiagnosisResponse | null;
  loading: boolean;
  error: string | null;
}

/**
 * useCoverage — loads FGOS coverage and gap diagnosis for a learner.
 * `lagModuleId` defaults to null (no gap diagnosis until requested).
 */
export function useCoverage(learnerId: string, lagModuleId?: string): GapCoverageState {
  const [state, setState] = useState<GapCoverageState>({
    coverage: null,
    gaps: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setState((prev) => ({ ...prev, loading: true, error: null }));
      try {
        const coverage = await getFgosCoverage(learnerId);
        const gaps = lagModuleId ? await diagnoseGaps(learnerId, lagModuleId) : null;
        if (!cancelled) setState({ coverage, gaps, loading: false, error: null });
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        if (!cancelled) setState((prev) => ({ ...prev, loading: false, error: message }));
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, [learnerId, lagModuleId]);

  return state;
}
