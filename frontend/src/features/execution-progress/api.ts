import { api } from '../../shared/api/client';
import type {
  CoverageResponse,
  GapDiagnosisResponse,
  Resource,
  ResourceListResponse,
} from '../../shared/api/types';

/** Fetch FGOS coverage for a learner. */
export async function getFgosCoverage(learnerId: string): Promise<CoverageResponse> {
  return api.get<CoverageResponse>(`/learners/${encodeURIComponent(learnerId)}/coverage/fgos`);
}

/** Diagnose gaps for a learner lag module. */
export async function diagnoseGaps(
  learnerId: string,
  lagModuleId: string,
): Promise<GapDiagnosisResponse> {
  return api.get<GapDiagnosisResponse>(
    `/learners/${encodeURIComponent(learnerId)}/gaps?lag_module_id=${encodeURIComponent(lagModuleId)}`,
  );
}

/** List resource catalog with optional filters. */
export async function listResources(
  params: { type?: string; format?: string } = {},
): Promise<ResourceListResponse> {
  const query = new URLSearchParams();
  if (params.type) query.set('type', params.type);
  if (params.format) query.set('format', params.format);
  return api.get<ResourceListResponse>(
    `/resources${query.toString() ? `?${query.toString()}` : ''}`,
  );
}

export type { Resource };
