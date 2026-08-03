/**
 * Practice Life — StoryCard and ProjectCard components.
 *
 * Stories are short contextual materials (≤5 min reading) linked to 1–3 topics.
 * Project ideas are cross-subject activities requiring modules from ≥2 subjects.
 */
import { Badge } from '../../shared/components';

export interface StoryCardData {
  id: string;
  title: string;
  linkedModules: string[];
  subjects: string[];
  realWorld: string;
  readingMinutes: number;
}

export interface ProjectCardData {
  id: string;
  title: string;
  modules: string[];
  difficultyLevel: string;
  expectedOutcome: string;
}

const difficultyColors: Record<string, 'success' | 'warning' | 'danger' | 'info'> = {
  basic: 'info',
  medium: 'warning',
  advanced: 'danger',
};

export function StoryCard({ story }: { story: StoryCardData }) {
  return (
    <div className="rounded-lg border border-neutral-200 bg-white p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-2">
        <h4 className="text-sm font-semibold text-neutral-800">{story.title}</h4>
        <span className="text-[10px] text-neutral-400 whitespace-nowrap">
          {story.readingMinutes} min read
        </span>
      </div>
      <div className="flex flex-wrap gap-1 mb-2">
        {story.subjects.map((s) => (
          <Badge key={s} color="info">
            {s}
          </Badge>
        ))}
      </div>
      <p className="text-xs text-neutral-600 mb-2 line-clamp-3">
        <strong>Real-world application:</strong> {story.realWorld}
      </p>
      <div className="text-[10px] text-neutral-400">
        Linked modules: {story.linkedModules.join(', ')}
      </div>
    </div>
  );
}

export function ProjectCard({ project }: { project: ProjectCardData }) {
  return (
    <div className="rounded-lg border border-green-200 bg-green-50 p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-2">
        <h4 className="text-sm font-semibold text-neutral-800">{project.title}</h4>
        <Badge color={difficultyColors[project.difficultyLevel] ?? 'info'}>
          {project.difficultyLevel}
        </Badge>
      </div>
      <p className="text-xs text-neutral-600 mb-2">{project.expectedOutcome}</p>
      <div className="text-[10px] text-neutral-400">Requires: {project.modules.join(', ')}</div>
    </div>
  );
}

export interface RecommendationPanelProps {
  moduleName: string;
  stories: StoryCardData[];
  projects: ProjectCardData[];
}

export function RecommendationPanel({ moduleName, stories, projects }: RecommendationPanelProps) {
  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold text-neutral-800">
        You've mastered <span className="text-primary-500">{moduleName}</span> — explore connections
      </h3>

      {stories.length > 0 && (
        <div className="space-y-3">
          <h4 className="text-sm font-medium text-neutral-700">
            📖 Related Stories ({stories.length})
          </h4>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {stories.map((s) => (
              <StoryCard key={s.id} story={s} />
            ))}
          </div>
        </div>
      )}

      {projects.length > 0 && (
        <div className="space-y-3">
          <h4 className="text-sm font-medium text-neutral-700">
            🔬 Project Ideas ({projects.length})
          </h4>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {projects.map((p) => (
              <ProjectCard key={p.id} project={p} />
            ))}
          </div>
        </div>
      )}

      {stories.length === 0 && projects.length === 0 && (
        <p className="text-sm text-neutral-500">
          No recommendations available for this module yet.
        </p>
      )}
    </div>
  );
}
