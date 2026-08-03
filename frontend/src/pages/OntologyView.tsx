import { OntologyBrowser, useOntology } from '../features/ontology-port';
import { RoleGate } from '../shared/guards/RoleGate';

/** OntologyView — knowledge-graph browser (protected, learner+). */
export function OntologyView() {
  const ontology = useOntology('math-5-1');
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'admin']}>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-gray-900">Ontology browser</h1>
        <OntologyBrowser ontology={ontology} />
      </div>
    </RoleGate>
  );
}
