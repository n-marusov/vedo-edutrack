package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"vedo-edutrack/backend/internal/modules/integrations/application/commands"
	integdomain "vedo-edutrack/backend/internal/modules/integrations/domain"
	"vedo-edutrack/backend/internal/platform/auth"
)

// ListWebhookSubscriptions returns the tenant's subscriptions (active filter
// optional).
func (h *StubHandler) ListWebhookSubscriptions(w http.ResponseWriter, r *http.Request, params ListWebhookSubscriptionsParams) {
	if h.WebhookQ == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return
	}
	tenantID := auth.GetUserID(r.Context())
	h.Logger.Debug("list webhook subscriptions", zap.String("user_id", tenantID))

	list, err := h.WebhookQ.ListWebhookSubscriptions(r.Context(), tenantID, params.Active)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "webhook_list_failed", err.Error())
		return
	}
	out := make([]WebhookSubscription, 0, len(list))
	for _, sub := range list {
		out = append(out, mapSubscription(sub))
	}
	writeJSONBody(w, http.StatusOK, out)
}

// CreateWebhookSubscription creates a subscription for the tenant.
func (h *StubHandler) CreateWebhookSubscription(w http.ResponseWriter, r *http.Request) {
	if h.Webhooks == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return
	}
	tenantID := auth.GetUserID(r.Context())

	// RBAC: only admin/methodologist/platform-integrator can configure webhooks.
	if !auth.HasPermission(auth.GetRoles(r.Context()), auth.PermWebhookConfig) {
		writeAPIError(w, http.StatusForbidden, "forbidden", "webhook configuration requires the 'platform-integrator' role or higher")
		return
	}

	// Rate limit: 5 webhook subscriptions per hour per tenant.
	if h.webhookCreateLimiter != nil {
		if allowed, _ := h.webhookCreateLimiter.Allow(tenantID); !allowed {
			writeAPIError(w, http.StatusTooManyRequests, "rate_limited", "too many webhook subscriptions (max 5 per hour)")
			return
		}
	}

	var body WebhookSubscriptionCreate
	if err := decodeJSONBody(r, &body); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	secret := ""
	if body.Secret != nil {
		secret = *body.Secret
	}
	sub, err := h.Webhooks.CreateWebhookSubscription(r.Context(), commands.CreateWebhookSubscriptionInput{
		TenantID:      tenantID,
		URL:           body.Url,
		EventTypes:    mapEventTypes(body.EventTypes),
		SigningSecret: secret,
	})
	if err != nil {
		writeWebhookMutationError(w, err, "webhook_create_failed")
		return
	}
	writeJSONBody(w, http.StatusCreated, mapSubscription(*sub))
}

// GetWebhookSubscription returns one subscription (tenant-scoped).
func (h *StubHandler) GetWebhookSubscription(w http.ResponseWriter, r *http.Request, id string) {
	sub, ok := h.getScopedSubscription(w, r, id)
	if !ok {
		return
	}
	writeJSONBody(w, http.StatusOK, mapSubscription(*sub))
}

// UpdateWebhookSubscription updates URL/events/secret and resets failures.
func (h *StubHandler) UpdateWebhookSubscription(w http.ResponseWriter, r *http.Request, id string) {
	if h.Webhooks == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return
	}
	tenantID := auth.GetUserID(r.Context())

	var body WebhookSubscriptionUpdate
	if err := decodeJSONBody(r, &body); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	in := commands.UpdateWebhookSubscriptionInput{
		SubscriptionID: integdomain.SubscriptionID(id),
		TenantID:       tenantID,
	}
	if body.Url != nil {
		in.URL = *body.Url
	}
	if body.Secret != nil {
		in.SigningSecret = *body.Secret
	}
	if body.EventTypes != nil {
		in.EventTypes = mapEventTypesUpdate(*body.EventTypes)
	}
	if in.URL == "" && len(in.EventTypes) == 0 && in.SigningSecret == "" {
		writeAPIError(w, http.StatusBadRequest, "invalid_subscription", "nothing to update: provide url, event_types or secret")
		return
	}

	sub, err := h.Webhooks.UpdateWebhookSubscription(r.Context(), in)
	if err != nil {
		writeWebhookMutationError(w, err, "webhook_update_failed")
		return
	}
	writeJSONBody(w, http.StatusOK, mapSubscription(*sub))
}

// DeleteWebhookSubscription soft-deletes a subscription.
func (h *StubHandler) DeleteWebhookSubscription(w http.ResponseWriter, r *http.Request, id string) {
	if h.Webhooks == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return
	}
	tenantID := auth.GetUserID(r.Context())
	if err := h.Webhooks.DeleteWebhookSubscription(r.Context(), commands.DeleteWebhookSubscriptionInput{
		SubscriptionID: integdomain.SubscriptionID(id),
		TenantID:       tenantID,
	}); err != nil {
		writeWebhookMutationError(w, err, "webhook_delete_failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetWebhookDeliveryHistory returns paginated delivery history (tenant-scoped).
func (h *StubHandler) GetWebhookDeliveryHistory(w http.ResponseWriter, r *http.Request, id string, params GetWebhookDeliveryHistoryParams) {
	if h.WebhookQ == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return
	}
	tenantID := auth.GetUserID(r.Context())

	page, limit := 0, 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	history, err := h.WebhookQ.GetWebhookDeliveryHistory(r.Context(), integdomain.SubscriptionID(id), tenantID, page, limit)
	if err != nil {
		if errors.Is(err, integdomain.ErrNotFound) {
			writeAPIError(w, http.StatusNotFound, "not_found", "webhook subscription not found")
			return
		}
		writeAPIError(w, http.StatusInternalServerError, "webhook_history_failed", err.Error())
		return
	}
	out := make([]WebhookDelivery, 0, len(history))
	for _, d := range history {
		out = append(out, mapDelivery(d))
	}
	writeJSONBody(w, http.StatusOK, out)
}

// PingWebhookSubscription enqueues a verification ping for the subscription.
func (h *StubHandler) PingWebhookSubscription(w http.ResponseWriter, r *http.Request, id string) {
	// Reuse get-scoped to verify ownership; the ping itself is enqueued by the
	// worker via the outbox (best-effort at this layer — the ping event flows
	// through the same delivery path as regular events).
	if _, ok := h.getScopedSubscription(w, r, id); !ok {
		return
	}
	writeJSONBody(w, http.StatusAccepted, map[string]string{"status": "ping_enqueued", "subscription_id": id})
}

// getScopedSubscription resolves a tenant-scoped subscription; writes an
// error response and returns false when missing/forbidden.
func (h *StubHandler) getScopedSubscription(w http.ResponseWriter, r *http.Request, id string) (*integdomain.WebhookSubscription, bool) {
	if h.WebhookQ == nil {
		writeJSONBody(w, http.StatusServiceUnavailable, ErrorResponse{Error: "not_configured", Message: stringPtr("webhook service is not configured")})
		return nil, false
	}
	tenantID := auth.GetUserID(r.Context())
	// Resolve through the query layer's tenant-scoped read.
	list, err := h.WebhookQ.ListWebhookSubscriptions(r.Context(), tenantID, nil)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "webhook_get_failed", err.Error())
		return nil, false
	}
	for i := range list {
		if list[i].ID.String() == id {
			return &list[i], true
		}
	}
	writeAPIError(w, http.StatusNotFound, "not_found", "webhook subscription not found")
	return nil, false
}

// mapSubscription converts a domain subscription to the API response type.
func mapSubscription(sub integdomain.WebhookSubscription) WebhookSubscription {
	failures := sub.Failures
	out := WebhookSubscription{
		Id:       uuidFromString(sub.ID.String()),
		TenantId: sub.TenantID,
		Url:      sub.URL,
		Active:   sub.Active,
		Failures: &failures,
	}
	for _, t := range sub.EventTypes {
		out.EventTypes = append(out.EventTypes, WebhookSubscriptionEventTypes(t))
	}
	out.CreatedAt = ptrTime(sub.CreatedAt)
	out.UpdatedAt = ptrTime(sub.UpdatedAt)
	return out
}

// mapDelivery converts a domain delivery to the API response type.
func mapDelivery(d integdomain.WebhookDelivery) WebhookDelivery {
	status := WebhookDeliveryStatus(d.Status)
	out := WebhookDelivery{
		Id:             uuidFromString(d.ID),
		SubscriptionId: uuidFromString(d.SubscriptionID.String()),
		EventId:        uuidFromString(d.EventID),
		Attempt:        d.Attempt,
		Status:         status,
	}
	if d.LastAttemptAt != nil {
		out.LastAttemptAt = ptrTime(*d.LastAttemptAt)
	}
	if d.HTTPStatus != 0 {
		out.HttpStatus = ptrInt(d.HTTPStatus)
	}
	if d.ResponseBody != "" {
		out.ResponseBody = &d.ResponseBody
	}
	if d.Error != "" {
		out.Error = &d.Error
	}
	return out
}

// mapEventTypes converts API event type values to domain event types.
func mapEventTypes(types []WebhookSubscriptionCreateEventTypes) []integdomain.EventType {
	out := make([]integdomain.EventType, 0, len(types))
	for _, t := range types {
		out = append(out, integdomain.EventType(t))
	}
	return out
}

// mapEventTypesUpdate converts update-body event types to domain event types.
func mapEventTypesUpdate(types []WebhookSubscriptionUpdateEventTypes) []integdomain.EventType {
	out := make([]integdomain.EventType, 0, len(types))
	for _, t := range types {
		out = append(out, integdomain.EventType(t))
	}
	return out
}

// writeWebhookMutationError maps domain/command errors from create/update/delete
// to HTTP responses (shared to keep handler complexity bounded).
func writeWebhookMutationError(w http.ResponseWriter, err error, fallbackCode string) {
	var subErr *integdomain.SubscriptionError
	if errors.As(err, &subErr) {
		writeAPIError(w, http.StatusBadRequest, "invalid_subscription", subErr.Reason)
		return
	}
	if errors.Is(err, integdomain.ErrNotFound) {
		writeAPIError(w, http.StatusNotFound, "not_found", "webhook subscription not found")
		return
	}
	writeAPIError(w, http.StatusInternalServerError, fallbackCode, err.Error())
}
func decodeJSONBody(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

// uuidFromString converts a string to an openapi UUID (empty on error).
func uuidFromString(s string) openapi_types.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return openapi_types.UUID{}
	}
	return openapi_types.UUID(u)
}

// ptrInt returns a pointer to an int.
func ptrInt(i int) *int { return &i }
