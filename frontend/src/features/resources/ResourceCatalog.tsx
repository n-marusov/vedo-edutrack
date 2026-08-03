import { useEffect, useState } from 'react';
import type { Resource } from '../../shared/api/types';
import { Badge } from '../../shared/components/Badge';
import { Card } from '../../shared/components/Card';
import { Input } from '../../shared/components/Input';
import { LoadingSpinner } from '../../shared/components/LoadingSpinner';
import { listResources } from '../execution-progress/api';

/** ResourceCatalog — catalog with format/source/difficulty filters (F3). */
export function ResourceCatalog() {
  const [resources, setResources] = useState<Resource[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [format, setFormat] = useState('');
  const [type, setType] = useState('');

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setLoading(true);
      setError(null);
      try {
        const response = await listResources({
          format: format || undefined,
          type: type || undefined,
        });
        if (!cancelled) setResources(response.items);
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        if (!cancelled) setError(message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
  }, [format, type]);

  return (
    <Card title="Resource catalog">
      <div className="mb-4 flex gap-3">
        <Input
          value={format}
          onChange={(event) => setFormat(event.target.value)}
          placeholder="Format filter (video, text, …)"
          aria-label="Format filter"
        />
        <Input
          value={type}
          onChange={(event) => setType(event.target.value)}
          placeholder="Type filter (content, enabling)"
          aria-label="Type filter"
        />
      </div>
      {loading && <LoadingSpinner label="Loading catalog…" />}
      {error && (
        <p role="alert" className="text-sm text-red-700">
          {error}
        </p>
      )}
      <ul data-testid="resource-list" className="space-y-2">
        {resources.map((resource) => (
          <li key={resource.id} className="flex items-center gap-3 text-sm text-gray-800">
            <span className="flex-1">{resource.title ?? resource.id}</span>
            {resource.format && <Badge color="info">{resource.format}</Badge>}
            {resource.type === 'enabling' && <Badge color="warning">enabling</Badge>}
          </li>
        ))}
      </ul>
    </Card>
  );
}
