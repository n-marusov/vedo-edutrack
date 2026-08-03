/**
 * Methodologist Dashboard — school-level aggregation per M2 specs.
 *
 * Widgets:
 * 1. FGOS coverage aggregated by class/school
 * 2. Top lagging topics (by number of learners with unmastered module)
 * 3. Ontology contribution (audit-log based)
 *
 * Role scope: methodologist sees only own school/classes.
 */
import { useMemo } from 'react';
import { Badge, Card } from '../../shared/components';

export interface ClassCoverage {
  className: string;
  learnerCount: number;
  fgosCoverage: number;
}

export interface LaggingTopic {
  moduleId: string;
  title: string;
  learnersWithGap: number;
  subject: string;
}

export interface MethodologistDashboardData {
  schoolName: string;
  classes: ClassCoverage[];
  laggingTopics: LaggingTopic[];
  ontologyContributions: number;
}

export interface MethodologistDashboardProps {
  data: MethodologistDashboardData;
}

export function MethodologistDashboard({ data }: MethodologistDashboardProps) {
  const schoolCoverage = useMemo(() => {
    if (data.classes.length === 0) {
      return 0;
    }
    const totalLearners = data.classes.reduce((sum, c) => sum + c.learnerCount, 0);
    if (totalLearners === 0) {
      return 0;
    }
    const weighted = data.classes.reduce((sum, c) => sum + c.fgosCoverage * c.learnerCount, 0);
    return Math.round(weighted / totalLearners);
  }, [data.classes]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-neutral-900">
          {data.schoolName} — Methodologist View
        </h2>
        <Badge color="info">{data.classes.length} classes</Badge>
      </div>

      {/* School-wide coverage */}
      <Card title="School FGOS Coverage">
        <div className="flex items-center gap-4">
          <div className="text-3xl font-bold text-neutral-800">{schoolCoverage}%</div>
          <div className="flex-1 h-3 bg-neutral-100 rounded-full overflow-hidden">
            <div
              className={`h-full rounded-full transition-all ${
                schoolCoverage >= 80
                  ? 'bg-green-500'
                  : schoolCoverage >= 50
                    ? 'bg-yellow-500'
                    : 'bg-red-500'
              }`}
              style={{ width: `${schoolCoverage}%` }}
            />
          </div>
        </div>
        <div className="mt-3 text-xs text-neutral-500">
          Weighted by learner count across classes
        </div>
      </Card>

      {/* Coverage by class */}
      <Card title="Coverage by Class">
        <div className="space-y-3">
          {data.classes.map((c) => (
            <div key={c.className} className="flex items-center gap-3">
              <span className="text-sm font-medium w-32 text-neutral-700">{c.className}</span>
              <div className="flex-1 h-2 bg-neutral-100 rounded-full overflow-hidden">
                <div
                  className={`h-full rounded-full ${
                    c.fgosCoverage >= 80
                      ? 'bg-green-500'
                      : c.fgosCoverage >= 50
                        ? 'bg-yellow-500'
                        : 'bg-red-500'
                  }`}
                  style={{ width: `${c.fgosCoverage}%` }}
                />
              </div>
              <span className="text-xs text-neutral-500 w-20 text-right">
                {c.fgosCoverage}% · {c.learnerCount} learners
              </span>
            </div>
          ))}
        </div>
      </Card>

      {/* Top lagging topics */}
      <Card title="Top Lagging Topics">
        {data.laggingTopics.length === 0 ? (
          <div className="text-sm text-neutral-500">No lagging topics detected.</div>
        ) : (
          <div className="space-y-2">
            {data.laggingTopics.map((topic) => (
              <div
                key={topic.moduleId}
                className="flex items-center justify-between rounded-lg bg-neutral-50 px-3 py-2"
              >
                <div>
                  <div className="text-sm font-medium text-neutral-800">{topic.title}</div>
                  <div className="text-xs text-neutral-500">
                    {topic.subject} · {topic.moduleId}
                  </div>
                </div>
                <Badge color={topic.learnersWithGap > 5 ? 'danger' : 'warning'}>
                  {topic.learnersWithGap} learners
                </Badge>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Ontology contribution */}
      <Card title="Ontology Contribution">
        <div className="flex items-center gap-3">
          <div className="text-2xl font-bold text-neutral-800">{data.ontologyContributions}</div>
          <div className="text-xs text-neutral-500">
            module bindings contributed this period (from audit log)
          </div>
        </div>
      </Card>
    </div>
  );
}
