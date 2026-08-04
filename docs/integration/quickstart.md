# Quickstart — VEDO EduTrack Integration (5 minutes)

This guide gets you from zero to a computed learning route with progress
tracking in about five minutes.

## 1. Get a token

Dev/staging issues JWT tokens directly:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"user-1","roles":["learner"]}' | jq -r .token)
```

> Production uses Keycloak (OIDC). See the Auth guide for the full flow.

## 2. Compute a route

```bash
curl -s -X POST http://localhost:8080/api/v1/routes/compute \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"learner_id":"user-1","goal_topic_id":"math-5-11"}' \
  | jq
```

The response contains the ordered route steps (strict prerequisites first),
horizons, and computed-at timestamp.

## 3. Check progress

```bash
curl -s http://localhost:8080/api/v1/learners/user-1/progress \
  -H "Authorization: Bearer $TOKEN" | jq
```

Shows plan-vs-actual per module: planned/actual dates, deviation days, and
readiness forecast.

## 4. Optional: receive webhooks

Create a subscription to get `module.mastered` events:

```bash
curl -s -X POST http://localhost:8080/api/v1/webhooks/subscriptions \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://your-app.example.com/hooks/edutrack","event_types":["module.mastered"]}'
```

See `webhook-guide.md` for signature verification and dedup handling.

## Next steps

- `sparql-guide.md` — query the ontology directly
- `mcp-guide.md` — connect AI agents to the EduTrack MCP server
- `examples/curl-examples.sh` — every endpoint in one script
