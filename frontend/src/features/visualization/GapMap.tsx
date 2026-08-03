/**
 * Gap Diagnostic Map — shows only root gaps (unmastered modules with no
 * unmastered prerequisites) with cascade arrows, ranked by impact.
 *
 * Reuses KnowledgeMap with root-gap filtering.
 */
import { useMemo } from 'react';
import { Badge, Card } from '../../shared/components';

export interface RootGap {
  moduleId: string;
  title: string;
  blockedModules: number;
  blockedSubjects: string[];
  rank: number;
}

export interface GapMapProps {
  gaps: RootGap[];
}

export function GapMap({ gaps }: GapMapProps) {
  const sorted = useMemo(() => [...gaps].sort((a, b) => a.rank - b.rank), [gaps]);

  if (sorted.length === 0) {
    return (
      <Card title="Gap Diagnostic Map">
        <div className="text-sm text-neutral-500">
          No root gaps detected. All prerequisites are mastered.
        </div>
      </Card>
    );
  }

  const topGap = sorted[0];

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold text-neutral-900">Gap Diagnostic Map</h2>

      {/* Summary */}
      <Card title="Root Gap Summary">
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <Badge color="danger">{sorted.length} root gaps</Badge>
            <span className="text-sm text-neutral-600">
              Top: close <strong>{topGap?.title}</strong> to unlock {topGap?.blockedModules} modules
              in {topGap?.blockedSubjects.join(', ')}
            </span>
          </div>
        </div>
      </Card>

      {/* Gap list */}
      <div className="space-y-3">
        {sorted.map((gap) => (
          <Card key={gap.moduleId}>
            <div className="flex items-start justify-between">
              <div>
                <div className="flex items-center gap-2">
                  <span className="text-xs text-neutral-400 font-mono">#{gap.rank}</span>
                  <span className="font-medium text-neutral-800">{gap.title}</span>
                  <Badge color="danger">Root Gap</Badge>
                </div>
                <div className="mt-1 text-xs text-neutral-500">
                  Blocks {gap.blockedModules} modules in {gap.blockedSubjects.length} subjects
                  {gap.blockedSubjects.length > 0 && <>: {gap.blockedSubjects.join(', ')}</>}
                </div>
              </div>
              <span className="text-xs font-mono text-neutral-400">{gap.moduleId}</span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}
