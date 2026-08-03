/**
 * Group Management Panel — mini-cards per learner with quick switching.
 *
 * Required fields per M2 spec:
 * - Name
 * - Current module
 * - FGOS coverage %
 * - Forecast status
 * - Attention flag (when lag > 10% OR unclosed prerequisites)
 */
import { useMemo } from 'react';
import { Badge, Card } from '../../shared/components';

export interface LearnerCard {
  id: string;
  name: string;
  currentModule: string;
  fgosCoverage: number;
  forecastStatus: 'on-track' | 'not-on-track' | 'at-risk';
  attentionFlag: boolean;
}

export interface GroupPanelProps {
  title: string;
  learners: LearnerCard[];
  onSelect?: (learnerId: string) => void;
}

export function GroupPanel({ title, learners, onSelect }: GroupPanelProps) {
  const atRisk = useMemo(() => learners.filter((l) => l.attentionFlag).length, [learners]);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-neutral-900">{title}</h2>
        <Badge color={atRisk > 0 ? 'warning' : 'info'}>
          {atRisk} of {learners.length} at risk
        </Badge>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        {learners.map((l) => (
          <Card
            key={l.id}
            className={l.attentionFlag ? 'border-l-4 border-l-yellow-500' : undefined}
          >
            <button
              type="button"
              className="w-full text-left"
              onClick={() => onSelect?.(l.id)}
              tabIndex={0}
              onKeyDown={(e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  onSelect?.(l.id);
                }
              }}
            >
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="font-medium text-neutral-800">{l.name}</span>
                  {l.attentionFlag && <Badge color="warning">⚠ Attention</Badge>}
                </div>
                <div className="text-xs text-neutral-500 space-y-0.5">
                  <div>Module: {l.currentModule}</div>
                  <div>
                    FGOS: {l.fgosCoverage}%
                    <span className="ml-1">
                      {l.fgosCoverage >= 80 ? '✅' : l.fgosCoverage >= 50 ? '⚠️' : '❌'}
                    </span>
                  </div>
                  <div>Forecast: {l.forecastStatus}</div>
                </div>
              </div>
            </button>
          </Card>
        ))}
      </div>
    </div>
  );
}
