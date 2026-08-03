/**
 * Custom knowledge node for React Flow.
 * Color-coded by mastery status.
 */
import { Handle, type NodeProps, Position } from '@xyflow/react';
import type { ModuleStatus } from './types';

const statusColors: Record<ModuleStatus, { bg: string; border: string; ring: string }> = {
  mastered: { bg: 'bg-green-100', border: 'border-green-500', ring: 'ring-green-300' },
  'in-progress': { bg: 'bg-yellow-50', border: 'border-yellow-500', ring: 'ring-yellow-300' },
  available: { bg: 'bg-blue-50', border: 'border-blue-400', ring: 'ring-blue-200' },
  blocked: { bg: 'bg-gray-100', border: 'border-gray-300', ring: 'ring-gray-200' },
  'unclosed-prereq': { bg: 'bg-red-50', border: 'border-red-500', ring: 'ring-red-300' },
};

export function KnowledgeNode({ data }: NodeProps) {
  const status: ModuleStatus = (data?.status as ModuleStatus) ?? 'available';
  const label = (data?.label as string) ?? '';
  const subject = (data?.subject as string) ?? '';
  const colors = statusColors[status];

  return (
    <div
      className={`px-3 py-2 rounded-lg border-2 shadow-sm ${colors.bg} ${colors.border} min-w-[120px]`}
    >
      <Handle type="target" position={Position.Top} className="!bg-neutral-400" />
      <div className="text-xs font-medium text-neutral-800 truncate max-w-[150px]">{label}</div>
      {subject && <div className="text-[10px] text-neutral-500 mt-0.5">{subject}</div>}
      <Handle type="source" position={Position.Bottom} className="!bg-neutral-400" />
    </div>
  );
}
