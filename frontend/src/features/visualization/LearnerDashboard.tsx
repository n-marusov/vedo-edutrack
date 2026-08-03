/**
 * Learner Dashboard — 6 mandatory widgets per M2 specs.
 *
 * Widgets:
 * 1. Current position (module + color status)
 * 2. Three horizons (far/mid/near)
 * 3. Plan-vs-actual deviation
 * 4. Subject progress summary
 * 5. FGOS coverage percentage
 * 6. Recommended stories/projects
 */
import { useMemo } from 'react';
import { Badge, Card } from '../../shared/components';
import type { ModuleStatus } from './types';

export interface Widget {
  title: string;
  value?: string | number;
  detail?: string;
  badge?: { text: string; color: 'success' | 'warning' | 'danger' | 'info' };
}

export interface LearnerDashboardData {
  learnerName: string;
  currentModule: string;
  currentStatus: ModuleStatus;
  horizons: { far: string; mid: string; near: string };
  planVsActual: { onTime: number; delayed: number; ahead: number };
  subjectProgress: Array<{ subject: string; percent: number }>;
  fgosCoverage: number;
  recommendedCount: number;
  lastUpdated: string;
}

export interface LearnerDashboardProps {
  data: LearnerDashboardData;
}

const statusLabels: Record<ModuleStatus, string> = {
  mastered: 'Mastered',
  'in-progress': 'In Progress',
  available: 'Available',
  blocked: 'Blocked',
  'unclosed-prereq': 'Gap',
};

const statusColors: Record<ModuleStatus, 'success' | 'warning' | 'danger' | 'info'> = {
  mastered: 'success',
  'in-progress': 'warning',
  available: 'info',
  blocked: 'info',
  'unclosed-prereq': 'danger',
};

export function LearnerDashboard({ data }: LearnerDashboardProps) {
  const widgets: Widget[] = useMemo(
    () => [
      {
        title: 'Current Position',
        value: data.currentModule,
        detail: `Status: ${statusLabels[data.currentStatus]}`,
        badge: { text: statusLabels[data.currentStatus], color: statusColors[data.currentStatus] },
      },
      {
        title: 'Three Horizons',
        value: `Far: ${data.horizons.far}`,
        detail: `Mid: ${data.horizons.mid}\nNear: ${data.horizons.near}`,
      },
      {
        title: 'Plan vs Actual',
        value: `${data.planVsActual.onTime} on time`,
        detail: `${data.planVsActual.delayed} delayed, ${data.planVsActual.ahead} ahead`,
        badge: data.planVsActual.delayed > 2 ? { text: 'Attention', color: 'warning' } : undefined,
      },
      {
        title: 'FGOS Coverage',
        value: `${data.fgosCoverage}%`,
        detail: `${100 - data.fgosCoverage}% remaining`,
        badge:
          data.fgosCoverage >= 80
            ? { text: 'Good', color: 'success' }
            : data.fgosCoverage >= 50
              ? { text: 'In Progress', color: 'warning' }
              : { text: 'Low', color: 'danger' },
      },
      {
        title: 'Subject Progress',
        value:
          data.subjectProgress.length > 0 && data.subjectProgress[0]
            ? `${data.subjectProgress[0].subject}: ${data.subjectProgress[0].percent}%`
            : 'No data',
        detail: data.subjectProgress
          .slice(1, 3)
          .map((s) => `${s.subject}: ${s.percent}%`)
          .join(', '),
      },
      {
        title: 'Recommendations',
        value: `${data.recommendedCount} items`,
        detail: 'Stories and project ideas',
      },
    ],
    [data],
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-neutral-900">
          {data.learnerName}&apos;s Dashboard
        </h2>
        <span className="text-xs text-neutral-400">Updated: {data.lastUpdated}</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {widgets.map((w) => (
          <Card key={w.title} title={w.title}>
            <div className="space-y-1">
              {w.badge && <Badge color={w.badge.color}>{w.badge.text}</Badge>}
              {w.value && <div className="text-lg font-semibold text-neutral-800">{w.value}</div>}
              {w.detail && (
                <div className="text-xs text-neutral-500 whitespace-pre-line">{w.detail}</div>
              )}
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
