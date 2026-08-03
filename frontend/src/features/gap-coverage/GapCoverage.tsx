import { Badge } from '../../shared/components/Badge';
import { Card } from '../../shared/components/Card';
import { LoadingSpinner } from '../../shared/components/LoadingSpinner';
import { useCoverage } from './useCoverage';

interface GapCoverageProps {
  learnerId: string;
  lagModuleId?: string;
}

/**
 * GapCoverage — root-cause gap viewer + live FGOS coverage with deficit list
 * and attestation readiness (F2).
 */
export function GapCoverage({ learnerId, lagModuleId }: GapCoverageProps) {
  const { coverage, gaps, loading, error } = useCoverage(learnerId, lagModuleId);

  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <Card title="FGOS coverage">
        {loading && <LoadingSpinner label="Loading coverage…" />}
        {error && (
          <p role="alert" className="text-sm text-red-700">
            {error}
          </p>
        )}
        {coverage && (
          <div data-testid="coverage-report">
            <div className="mb-4 h-3 w-full overflow-hidden rounded-full bg-gray-100">
              <div
                className="h-full rounded-full bg-blue-600"
                style={{ width: `${Math.min(coverage.percent, 100)}%` }}
              />
            </div>
            <p className="text-sm text-gray-700">
              {coverage.percent.toFixed(1)}% covered ({coverage.covered}/{coverage.total})
            </p>
            <div className="mt-2">
              {coverage.ready ? (
                <Badge color="success">attestation ready</Badge>
              ) : (
                <Badge color="warning">not ready</Badge>
              )}
            </div>
            {(coverage.deficits ?? []).length > 0 && (
              <ul className="mt-3 space-y-1 text-sm text-gray-600">
                {coverage.deficits?.map((deficit) => (
                  <li key={deficit.requirement_id}>
                    {deficit.requirement_id}
                    {deficit.blocking_module_id
                      ? ` (blocked by ${deficit.blocking_module_id})`
                      : ''}
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </Card>

      <Card title="Gap diagnosis">
        {lagModuleId ? (
          gaps && (
            <div data-testid="gap-diagnosis">
              {gaps.root_causes.length === 0 ? (
                <p className="text-sm text-gray-600">No root cause found.</p>
              ) : (
                <ul className="space-y-2">
                  {gaps.root_causes.map((cause) => (
                    <li
                      key={cause.module_id}
                      className="flex items-center gap-3 text-sm text-gray-800"
                    >
                      <span className="flex-1">{cause.module_id}</span>
                      <Badge color="danger">root</Badge>
                      {cause.blocked_modules !== undefined && (
                        <span className="text-xs text-gray-500">
                          blocks {cause.blocked_modules}
                        </span>
                      )}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )
        ) : (
          <p className="text-sm text-gray-500">Pass a lag module id to run diagnosis.</p>
        )}
      </Card>
    </div>
  );
}
