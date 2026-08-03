import { Background, Controls, type Edge, MiniMap, type Node, ReactFlow } from '@xyflow/react';
/**
 * Visualization feature — Knowledge Map component.
 *
 * Renders a 2D knowledge graph with color-coded nodes and two viewing modes:
 * - Critical path (strict + essential edges only)
 * - Exploration (all edges)
 *
 * Uses @xyflow/react (React Flow) for graph rendering.
 */
import { useMemo, useState } from 'react';
import '@xyflow/react/dist/style.css';
import { Badge } from '../../shared/components';
import { KnowledgeNode } from './KnowledgeNode';
import { Legend } from './Legend';
import type { GraphEdge, GraphMode, GraphNode, ModuleStatus } from './types';

export interface KnowledgeMapProps {
  nodes: GraphNode[];
  edges: GraphEdge[];
  progress: Record<string, ModuleStatus>;
}

function toFlowNode(n: GraphNode, progress: Record<string, ModuleStatus>): Node {
  const status = progress[n.id] ?? 'available';
  return {
    id: n.id,
    type: 'knowledge',
    position: { x: n.x ?? Math.random() * 600, y: n.y ?? Math.random() * 400 },
    data: { label: n.title, status, subject: n.subject },
  };
}

function toFlowEdge(e: GraphEdge): Edge {
  return {
    id: e.id,
    source: e.source,
    target: e.target,
    type: 'smoothstep',
    animated: e.type === 'hasStrictPrerequisite',
    style: { stroke: edgeColor(e.type), strokeWidth: e.type === 'hasStrictPrerequisite' ? 2 : 1 },
    data: { linkType: e.type },
  };
}

function edgeColor(linkType: string): string {
  switch (linkType) {
    case 'hasStrictPrerequisite':
      return '#ef4444'; // red-500
    case 'hasSoftPrerequisite':
      return '#eab308'; // yellow-500
    case 'enriches':
      return '#3b82f6'; // blue-500
    case 'appliesTo':
      return '#22c55e'; // green-500
    case 'isAnalogousTo':
      return '#a855f7'; // purple-500
    default:
      return '#9ca3af'; // gray-400
  }
}

export function KnowledgeMap({ nodes, edges, progress }: KnowledgeMapProps) {
  const [mode, setMode] = useState<GraphMode>('exploration');

  const filteredEdges = useMemo(() => {
    if (mode === 'critical-path') {
      return edges.filter((e) => e.type === 'hasStrictPrerequisite');
    }
    return edges;
  }, [edges, mode]);

  const flowNodes = useMemo(() => nodes.map((n) => toFlowNode(n, progress)), [nodes, progress]);
  const flowEdges = useMemo(() => filteredEdges.map(toFlowEdge), [filteredEdges]);

  const nodeTypes = useMemo(() => ({ knowledge: KnowledgeNode }), []);

  const stats = useMemo(() => {
    const counts: Record<ModuleStatus, number> = {
      mastered: 0,
      'in-progress': 0,
      available: 0,
      blocked: 0,
      'unclosed-prereq': 0,
    };
    for (const n of nodes) {
      const s = progress[n.id] ?? 'available';
      counts[s]++;
    }
    return counts;
  }, [nodes, progress]);

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => setMode('critical-path')}
            className={`px-3 py-1 text-sm rounded transition-colors ${
              mode === 'critical-path'
                ? 'bg-primary-500 text-white'
                : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'
            }`}
          >
            Critical Path
          </button>
          <button
            type="button"
            onClick={() => setMode('exploration')}
            className={`px-3 py-1 text-sm rounded transition-colors ${
              mode === 'exploration'
                ? 'bg-primary-500 text-white'
                : 'bg-neutral-100 text-neutral-700 hover:bg-neutral-200'
            }`}
          >
            Exploration
          </button>
        </div>
        <div className="flex gap-2 text-xs">
          <Badge color="success">{stats.mastered} mastered</Badge>
          <Badge color="warning">{stats['in-progress']} in progress</Badge>
          <Badge color="info">{stats.available} available</Badge>
          <Badge color="danger">{stats['unclosed-prereq']} gaps</Badge>
        </div>
      </div>

      {/* Graph */}
      <div className="border border-neutral-200 rounded-lg overflow-hidden" style={{ height: 500 }}>
        <ReactFlow
          nodes={flowNodes}
          edges={flowEdges}
          nodeTypes={nodeTypes}
          fitView
          fitViewOptions={{ padding: 0.2 }}
          attributionPosition="bottom-left"
        >
          <Background />
          <Controls />
          <MiniMap
            nodeColor={(n) => {
              switch (n.data?.status) {
                case 'mastered':
                  return '#22c55e';
                case 'in-progress':
                  return '#eab308';
                case 'unclosed-prereq':
                  return '#ef4444';
                default:
                  return '#3b82f6';
              }
            }}
          />
        </ReactFlow>
      </div>
      <Legend />
    </div>
  );
}
