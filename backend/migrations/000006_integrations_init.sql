-- create schema "integrations"
CREATE SCHEMA IF NOT EXISTS integrations;

-- create "webhook_subscriptions" (subscriber endpoints bound to event types)
-- Validation (REQ-FR-api.webhooks.*):
--   event_types ⊆ {module.mastered, plan.deviated, route.recalculated, standard.risk_detected}
--   url is https (http allowed for localhost dev sandbox)
--   active subscriptions are capped at MaxSubscriptionsPerTenant (10)
--   failures counts consecutive delivery failures; >= 5 deactivates the subscription
CREATE TABLE IF NOT EXISTS integrations.webhook_subscriptions (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    tenant_id text NOT NULL,
    url text NOT NULL,
    event_types text[] NOT NULL,
    signing_secret text NOT NULL,
    active boolean NOT NULL DEFAULT true,
    failures integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT webhook_subscriptions_failures_check CHECK (failures >= 0),
    CONSTRAINT webhook_subscriptions_active_check CHECK (active IN (true, false))
);

-- create "webhook_deliveries" (delivery record per outbox event per subscription)
-- Status lifecycle: pending -> sent | failed -> permanent_failure.
CREATE TABLE IF NOT EXISTS integrations.webhook_deliveries (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    subscription_id uuid NOT NULL,
    event_id uuid NOT NULL,
    attempt integer NOT NULL DEFAULT 1,
    status text NOT NULL DEFAULT 'pending',
    last_attempt_at timestamptz,
    http_status integer,
    response_body text NOT NULL DEFAULT '',
    error text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT webhook_deliveries_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES integrations.webhook_subscriptions (id) ON DELETE CASCADE,
    CONSTRAINT webhook_deliveries_status_check CHECK (status IN ('pending', 'sent', 'failed', 'permanent_failure')),
    CONSTRAINT webhook_deliveries_attempt_check CHECK (attempt >= 1),
    CONSTRAINT webhook_deliveries_unique_delivery UNIQUE (subscription_id, event_id)
);

-- create "event_dedup" (idempotent webhook receive: (source, event_id) unique)
CREATE TABLE IF NOT EXISTS integrations.event_dedup (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    source text NOT NULL,
    event_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT event_dedup_source_event_key UNIQUE (source, event_id)
);

-- create "outbox_events" (transactional outbox for outbound webhook events)
-- Dequeue order: pending first, ordered by created_at; retry_count drives the
-- exponential backoff schedule (1s, 2s, 4s, 8s, 16s — max 5 attempts).
CREATE TABLE IF NOT EXISTS integrations.outbox_events (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL,
    event_type text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    status text NOT NULL DEFAULT 'pending',
    retry_count integer NOT NULL DEFAULT 0,
    delivered_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT outbox_events_event_id_key UNIQUE (event_id),
    CONSTRAINT outbox_events_status_check CHECK (status IN ('pending', 'in_flight', 'delivered', 'failed'))
);

-- create indexes for worker polling and delivery history lookups
CREATE INDEX IF NOT EXISTS outbox_events_status_created_idx ON integrations.outbox_events (status, created_at);
CREATE INDEX IF NOT EXISTS webhook_subscriptions_tenant_id_idx ON integrations.webhook_subscriptions (tenant_id);
CREATE INDEX IF NOT EXISTS webhook_deliveries_subscription_idx ON integrations.webhook_deliveries (subscription_id);
CREATE INDEX IF NOT EXISTS webhook_deliveries_event_idx ON integrations.webhook_deliveries (event_id);

-- down: DROP TABLE IF EXISTS integrations.outbox_events;
-- down: DROP TABLE IF EXISTS integrations.event_dedup;
-- down: DROP TABLE IF EXISTS integrations.webhook_deliveries;
-- down: DROP TABLE IF EXISTS integrations.webhook_subscriptions;
-- down: DROP SCHEMA IF EXISTS integrations;
