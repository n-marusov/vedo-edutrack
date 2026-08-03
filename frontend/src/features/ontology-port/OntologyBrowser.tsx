import { useMemo } from 'react';
import { Badge } from '../../shared/components/Badge';
import { Card } from '../../shared/components/Card';
import { LoadingSpinner } from '../../shared/components/LoadingSpinner';
import type { GraphNode } from './useOntology';
import type { OntologyBrowserState } from './useOntology';

/** Edge colors by link type (F0 five link types). */
const edgeColors: Record<string, string> = {
  hasStrictPrerequisite: '#dc2626', // red
  hasSoftPrerequisite: '#f59e0b', // amber
  enriches: '#10b981', // emerald
  appliesTo: '#3b82f6', // blue
  isAnalogousTo: '#8b5cf6', // violet
};

const linkTypeLabel: Record<string, string> = {
  hasStrictPrerequisite: 'strict',
  hasSoftPrerequisite: 'soft',
  enriches: 'enriches',
  appliesTo: 'appliesTo',
  isAnalogousTo: 'analogous',
};

interface OntologyBrowserProps {
  ontology: OntologyBrowserState;
}

/**
 * OntologyBrowser — renders the ontology graph as an SVG with color-coded
 * link edges plus a module detail panel (F4 knowledge-map preview).
 */
export function OntologyBrowser({ ontology }: OntologyBrowserProps) {
  const layout = useMemo(() => layoutNodes(ontology.nodes), [ontology.nodes]);

  const selected = ontology.nodes.find((node) => node.id === ontology.selectedId) ?? null;

  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <Card title="Knowledge graph" className="lg:col-span-2">
        {ontology.error && (
          <p role="alert" className="mb-3 text-sm text-red-700">
            {ontology.error}
          </p>
        )}
        {ontology.loading && <LoadingSpinner label="Loading ontology…" />}
        <svg
          viewBox="0 0 640 480"
          className="h-96 w-full rounded border border-gray-100"
          data-testid="ontology-graph"
          role="img"
          aria-label="Knowledge graph"
        >
          {ontology.edges.map(([source, target, linkType]) => {
            const from = layout.get(source);
            const to = layout.get(target);
            if (!from || !to) return null;
            return (
              <line
                key={`edge-${source}-${target}`}
                x1={from.x}
                y1={from.y}
                x2={to.x}
                y2={to.y}
                stroke={edgeColors[linkType] ?? '#9ca3af'}
                strokeWidth={2}
                strokeDasharray={
                  linkType === 'hasSoftPrerequisite' || linkType === 'isAnalogousTo'
                    ? '4 3'
                    : undefined
                }
                data-testid={`edge-${source}-${target}`}
              />
            );
          })}
          {ontology.nodes.map((node) => {
            const position = layout.get(node.id);
            if (!position) return null;
            const isSelected = node.id === ontology.selectedId;
            return (
              <g
                key={node.id}
                transform={`translate(${position.x}, ${position.y})`}
                // biome-ignore lint/a11y/useSemanticElements: SVG graph node — a <button> inside <svg> is invalid; keyboard access is provided via onKeyDown.
                role="button"
                tabIndex={0}
                onClick={() => ontology.selectNode(node.id)}
                onKeyDown={(event) => {
                  if (event.key === 'Enter' || event.key === ' ') {
                    event.preventDefault();
                    ontology.selectNode(node.id);
                  }
                }}
                className="cursor-pointer"
                data-testid={`node-${node.id}`}
              >
                <circle
                  r={14}
                  fill={isSelected ? '#2563eb' : '#ffffff'}
                  stroke="#374151"
                  strokeWidth={2}
                />
                <text
                  textAnchor="middle"
                  dy={30}
                  className="select-none"
                  fontSize={10}
                  fill="#374151"
                >
                  {truncate(node.id, 14)}
                </text>
              </g>
            );
          })}
        </svg>
        <div className="mt-3 flex flex-wrap gap-2">
          {Object.entries(linkTypeLabel).map(([type, label]) => (
            <span key={type} className="inline-flex items-center gap-1.5 text-xs text-gray-600">
              <span
                className="inline-block h-2 w-4 rounded"
                style={{ backgroundColor: edgeColors[type] }}
              />
              {label}
            </span>
          ))}
        </div>
      </Card>

      <Card title="Module detail">
        {selected ? (
          <div data-testid="module-detail">
            <h4 className="text-sm font-semibold text-gray-900">{selected.title}</h4>
            {selected.linkType && (
              <div className="mt-2">
                <Badge color="info">{linkTypeLabel[selected.linkType] ?? selected.linkType}</Badge>
              </div>
            )}
            {selected.description ? (
              <p className="mt-3 text-sm text-gray-600">{selected.description}</p>
            ) : (
              <p className="mt-3 text-sm text-gray-400">No description.</p>
            )}
            <button
              type="button"
              onClick={() => void ontology.expandNode(selected.id)}
              className="mt-4 rounded bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            >
              Expand module
            </button>
          </div>
        ) : (
          <p className="text-sm text-gray-500">Select a module to inspect it.</p>
        )}
      </Card>
    </div>
  );
}

/** Simple circular layout with a small random-free deterministic offset. */
function layoutNodes(nodes: GraphNode[]): Map<string, { x: number; y: number }> {
  const layout = new Map<string, { x: number; y: number }>();
  const count = nodes.length;
  nodes.forEach((node, index) => {
    if (count === 1) {
      layout.set(node.id, { x: 320, y: 240 });
      return;
    }
    const angle = (2 * Math.PI * index) / count - Math.PI / 2;
    const radius = Math.min(200, 60 + 40 * Math.sqrt(count));
    layout.set(node.id, {
      x: 320 + radius * Math.cos(angle),
      y: 240 + radius * Math.sin(angle),
    });
  });
  return layout;
}

function truncate(value: string, max: number): string {
  return value.length > max ? `${value.slice(0, max - 1)}…` : value;
}
