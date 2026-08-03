import { useState } from 'react';
import { Badge } from '../../shared/components/Badge';
import { Button } from '../../shared/components/Button';
import { Card } from '../../shared/components/Card';
import { Input } from '../../shared/components/Input';
import { LoadingSpinner } from '../../shared/components/LoadingSpinner';
import type { Horizon, RouteBuilderState } from './useRouteCompute';
import { horizonSlice } from './useRouteCompute';

const horizonLabels: Record<Horizon, string> = {
  far: 'Far (full path)',
  mid: 'Mid (next modules)',
  near: 'Near (current module)',
};

interface RouteBuilderProps {
  route: RouteBuilderState;
}

/**
 * RouteBuilder — goal selector + route computation with three-horizon view
 * (F1: shortest path with strict/soft prerequisites).
 */
export function RouteBuilder({ route }: RouteBuilderProps) {
  const [goal, setGoal] = useState('');

  const visible = horizonSlice(route.route, route.horizon);

  return (
    <div className="space-y-4">
      <Card title="Goal selector">
        <div className="flex gap-3">
          <Input
            value={goal}
            onChange={(event) => setGoal(event.target.value)}
            placeholder="Goal module id, e.g. math-5-11"
            aria-label="Goal module id"
          />
          <Button onClick={() => void route.compute(goal)} disabled={route.loading}>
            Compute route
          </Button>
        </div>
        {route.error && (
          <p role="alert" className="mt-3 text-sm text-red-700">
            {route.error}
          </p>
        )}
      </Card>

      {route.loading && <LoadingSpinner label="Computing route…" />}

      {route.recomputeNeeded && (
        <Badge color="warning">Route needs recalculation (ontology/progress changed)</Badge>
      )}

      {route.route.length > 0 && (
        <Card
          title="Computed route"
          header={
            <div className="flex flex-wrap gap-2">
              {(Object.keys(horizonLabels) as Horizon[]).map((horizon) => (
                <button
                  key={horizon}
                  type="button"
                  onClick={() => route.setHorizon(horizon)}
                  className={`rounded px-2.5 py-1 text-xs font-medium ${
                    route.horizon === horizon
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-100 text-gray-700'
                  }`}
                  aria-pressed={route.horizon === horizon}
                >
                  {horizonLabels[horizon]}
                </button>
              ))}
            </div>
          }
        >
          <ol data-testid="route-timeline" className="space-y-2">
            {visible.map((topic) => (
              <li key={topic.topic_id} className="flex items-center gap-3 text-sm text-gray-800">
                <span className="inline-flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-800">
                  {topic.order + 1}
                </span>
                <span className="flex-1">{topic.topic_id}</span>
                {topic.link_type && <Badge color="info">{topic.link_type}</Badge>}
              </li>
            ))}
          </ol>
        </Card>
      )}
    </div>
  );
}
