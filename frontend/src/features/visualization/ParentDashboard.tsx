/**
 * Parent / HR Dashboard — 5 mandatory widgets per M2 specs.
 *
 * Widgets:
 * 1. Progress overview (child status summary)
 * 2. FGOS coverage percentage
 * 3. Deviations with magnitude + color highlight (lag > 10% = signal)
 * 4. Forecast to checkpoint
 * 5. Recommendations from lag areas
 *
 * Supports 2+ children switching.
 */
import { useMemo, useState } from 'react';
import { Badge, Card } from '../../shared/components';
import type { LearnerDashboardData } from './LearnerDashboard';
import { LearnerDashboard } from './LearnerDashboard';

export interface ChildOverview {
  id: string;
  name: string;
  fgosCoverage: number;
  forecastStatus: 'on-track' | 'not-on-track' | 'at-risk';
  lagPercent: number;
  lagModule: string;
  recommendedCount: number;
}

export interface ParentDashboardData {
  parentName: string;
  children: ChildOverview[];
  childDetail: Record<string, LearnerDashboardData>;
}

export interface ParentDashboardProps {
  data: ParentDashboardData;
}

export function ParentDashboard({ data }: ParentDashboardProps) {
  const [selectedChildId, setSelectedChildId] = useState(data.children[0]?.id);

  const selectedChild = useMemo(
    () => data.children.find((c) => c.id === selectedChildId),
    [data.children, selectedChildId],
  );
  const selectedDetail = selectedChild ? data.childDetail[selectedChild.id] : undefined;

  const atRisk = useMemo(
    () => data.children.filter((c) => c.forecastStatus !== 'on-track').length,
    [data.children],
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-neutral-900">
          {data.parentName}&apos;s Family Dashboard
        </h2>
        <Badge color={atRisk > 0 ? 'warning' : 'info'}>
          {atRisk} of {data.children.length} need attention
        </Badge>
      </div>

      {/* Child switcher */}
      {data.children.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {data.children.map((child) => (
            <button
              key={child.id}
              type="button"
              onClick={() => setSelectedChildId(child.id)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                child.id === selectedChildId
                  ? 'bg-primary-500 text-white'
                  : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'
              }`}
            >
              {child.name}
              {child.lagPercent > 10 && ' ⚠️'}
            </button>
          ))}
        </div>
      )}

      {selectedChild && selectedDetail ? (
        <>
          {/* Summary strip for the selected child */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
            <Card title="FGOS Coverage">
              <div className="text-lg font-semibold">{selectedChild.fgosCoverage}%</div>
              <div className="text-xs text-neutral-500">
                {selectedChild.fgosCoverage >= 80 ? 'on track for attestation' : 'gap remaining'}
              </div>
            </Card>
            <Card title="Forecast">
              <Badge
                color={
                  selectedChild.forecastStatus === 'on-track'
                    ? 'success'
                    : selectedChild.forecastStatus === 'at-risk'
                      ? 'warning'
                      : 'danger'
                }
              >
                {selectedChild.forecastStatus}
              </Badge>
            </Card>
            <Card title="Deviation">
              <div
                className={`text-lg font-semibold ${
                  selectedChild.lagPercent > 10 ? 'text-red-600' : 'text-green-600'
                }`}
              >
                {selectedChild.lagPercent}%
              </div>
              <div className="text-xs text-neutral-500">
                {selectedChild.lagPercent > 10
                  ? `⚠ lagging on ${selectedChild.lagModule}`
                  : 'within tolerance'}
              </div>
            </Card>
            <Card title="Recommendations">
              <div className="text-lg font-semibold">{selectedChild.recommendedCount}</div>
              <div className="text-xs text-neutral-500">stories & projects</div>
            </Card>
          </div>

          {/* Full learner dashboard for the selected child */}
          <LearnerDashboard data={selectedDetail} />
        </>
      ) : (
        <Card title="No children">
          <div className="text-sm text-neutral-500">
            No learners linked to this parent account yet.
          </div>
        </Card>
      )}
    </div>
  );
}
