import { api } from '../../shared/api/client';
import type { Concept, ConceptResponse } from '../../shared/api/types';

/** Fetch one ontology concept (with linked topics) from the REST API. */
export async function getConcept(topicId: string): Promise<Concept> {
  const response = await api.get<ConceptResponse>(
    `/ontology/concepts?topic_id=${encodeURIComponent(topicId)}`,
  );
  return response.concept;
}

/** Fetch resources bound to a module. */
export async function getModuleResources(moduleId: string) {
  return api.get(`/modules/${encodeURIComponent(moduleId)}/resources`);
}
