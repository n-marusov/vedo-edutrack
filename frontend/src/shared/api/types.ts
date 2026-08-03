/** Shared API types — mirror of the OpenAPI spec (backend/api/openapi/v1.yaml). */

export interface TokenRequest {
  user_id: string;
  roles: string[];
}

export interface TokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
}

export interface UserInfo {
  user_id: string;
  roles: string[];
}

export interface RouteComputeRequest {
  learner_id: string;
  goal_topic_id: string;
}

export interface RouterTopic {
  topic_id: string;
  order: number;
  link_type?: string;
}

export interface RouteComputeResponse {
  route: RouterTopic[];
}

export interface ConceptLink {
  topic_id: string;
  link_type: string;
}

export interface Concept {
  id: string;
  title: string;
  description?: string;
  links?: ConceptLink[];
}

export interface ConceptResponse {
  concept: Concept;
}

/** Resource catalog entry (M1, F3). */
export interface Resource {
  id: string;
  title?: string;
  type: 'content' | 'enabling';
  format?: string;
  source?: string;
  difficulty?: string;
  duration_minutes?: number;
  cost?: number;
  uri?: string;
}

export interface ResourceListResponse {
  items: Resource[];
  total: number;
}

/** Root-cause gap diagnosis (M1, F2). */
export interface RootCause {
  module_id: string;
  mastery?: number;
  blocked_modules?: number;
}

export interface GapDiagnosisResponse {
  status: 'root-causes-found' | 'no-root-cause-found';
  root_causes: RootCause[];
}

export interface Deficit {
  requirement_id: string;
  blocking_module_id?: string;
}

export interface CoverageResponse {
  covered: number;
  total: number;
  percent: number;
  ready?: boolean;
  deficits?: Deficit[];
}

export interface ApiErrorBody {
  error: string;
  message?: string;
  endpoint?: string;
}
