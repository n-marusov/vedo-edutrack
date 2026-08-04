#!/usr/bin/env bash
# VEDO EduTrack — curl examples for every public endpoint (dev/staging).
# Requires: curl, jq.
set -euo pipefail

API_BASE="${API_BASE:-http://localhost:8080/api/v1}"
TOKEN=${TOKEN:-}

if [[ -z "$TOKEN" ]]; then
  TOKEN=$(curl -s -X POST "$API_BASE/auth/token" \
    -H 'Content-Type: application/json' \
    -d '{"user_id":"user-1","roles":["learner"]}' | jq -r .token)
fi

AUTH=(-H "Authorization: Bearer $TOKEN")

echo "== Health =="
curl -s "$API_BASE/../healthz" | jq

echo "== Me =="
curl -s "${AUTH[@]}" "$API_BASE/me" | jq

echo "== Compute route =="
curl -s "${AUTH[@]}" -X POST "$API_BASE/routes/compute" \
  -H 'Content-Type: application/json' \
  -d '{"learner_id":"user-1","goal_topic_id":"math-5-11"}' | jq

echo "== Progress =="
curl -s "${AUTH[@]}" "$API_BASE/learners/user-1/progress" | jq

echo "== FGOS coverage =="
curl -s "${AUTH[@]}" "$API_BASE/learners/user-1/coverage/fgos" | jq

echo "== Gap diagnosis =="
curl -s "${AUTH[@]}" "$API_BASE/learners/user-1/gaps?lag_module_id=chemistry" | jq

echo "== SPARQL (list classes) =="
QUERY='SELECT ?s ?label WHERE { ?s a <https://www.w3.org/2002/07/owl#Class> }'
curl -s -G "${AUTH[@]}" "$API_BASE/sparql" --data-urlencode "query=$QUERY" | jq

echo "== Create webhook subscription =="
curl -s "${AUTH[@]}" -X POST "$API_BASE/webhooks/subscriptions" \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://your-app.example.com/hooks/edutrack","event_types":["module.mastered"],"secret":"01234567890123456789012345678901"}' | jq

echo "== List webhook subscriptions =="
curl -s "${AUTH[@]}" "$API_BASE/webhooks/subscriptions" | jq
