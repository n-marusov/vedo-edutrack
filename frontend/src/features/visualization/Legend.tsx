/**
 * Color legend for the knowledge map.
 */
const items = [
  { status: 'mastered', color: 'bg-green-500', label: 'Mastered' },
  { status: 'in-progress', color: 'bg-yellow-500', label: 'In Progress' },
  { status: 'available', color: 'bg-blue-500', label: 'Available' },
  { status: 'blocked', color: 'bg-gray-400', label: 'Blocked' },
  { status: 'unclosed-prereq', color: 'bg-red-500', label: 'Unclosed Prereq' },
] as const;

export function Legend() {
  return (
    <div className="flex flex-wrap gap-4 text-xs text-neutral-600">
      {items.map((item) => (
        <div key={item.status} className="flex items-center gap-1.5">
          <span className={`inline-block w-3 h-3 rounded ${item.color}`} />
          <span>{item.label}</span>
        </div>
      ))}
    </div>
  );
}
