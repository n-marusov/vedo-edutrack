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

export interface ApiErrorBody {
  error: string;
  message?: string;
  endpoint?: string;
}
