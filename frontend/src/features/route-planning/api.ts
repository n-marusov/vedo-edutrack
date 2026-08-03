import { api } from '../../shared/api/client';
import type { RouteComputeResponse, RouterTopic } from '../../shared/api/types';

/** Compute a learning route for a learner toward a goal topic. */
export async function computeRoute(learnerId: string, goalTopicId: string): Promise<RouterTopic[]> {
  const response = await api.post<RouteComputeResponse>('/routes/compute', {
    learner_id: learnerId,
    goal_topic_id: goalTopicId,
  });
  return response.route;
}
